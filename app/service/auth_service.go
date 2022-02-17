package service

import (
	"database/sql"
	"errors"
	"fiber-starter/app/api/requests"
	"fiber-starter/app/model"
	"fiber-starter/app/repository"
	"fiber-starter/config"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"fiber-starter/pkg/utils"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GenerateOTPToken(dbctx db.DBCtx, channel, email, phone string) (string, error)
	SendOTPToken(cfg *config.RabbitMQ, emailPhone, channel, usedFor, otpToken, via string) error
	Registration(dbctx db.DBCtx, req requests.RegisterRequest, channel string, rabbitCfg *config.RabbitMQ) (model.User, string, error)
	AuthLogin(dbctx db.DBCtx, req requests.LoginRequest, channel string, rabbitCfg *config.RabbitMQ) (model.User, string, error)
	SendOTPTokenByType(dbctx db.DBCtx, emailPhone, channel, sendType string, rabbitCfg *config.RabbitMQ) (string, bool, error)
	ChangePassword(dbctx db.DBCtx, req requests.ResetPasswordRequest, channel string) error
	OTPTokenValidation(dbctx db.DBCtx, req requests.ValidateOTPTokenRequest, channel, sendType string) error
}

type authService struct {
	userR repository.UserRepository
	roleR repository.RoleRepository
	otpR  repository.UserOTPRepository
}

func NewAuthService(user repository.UserRepository, role repository.RoleRepository, otp repository.UserOTPRepository) *authService {
	return &authService{user, role, otp}
}

func (s *authService) GenerateOTPToken(dbctx db.DBCtx, channel, email, phone string) (string, error) {
	// Generate otpToken channel web adm cms
	if channel == utils.ChannelWeb {
		token, _, err := utils.GenerateJWT(0, "", phone, email, "")
		return token, err
	}

	// Define data otp exist
	var otp string
	userOTP, err := s.otpR.GetByEmailOrPhone(dbctx, email, phone)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return otp, err
	}

	// Adjustment send otp in one row
	if userOTP.ID > 0 {
		userOTP, err = s.otpR.UpdateOTP(dbctx, userOTP)
	} else {
		userOTP, err = s.otpR.Insert(dbctx, model.UserOTP{
			Email: email,
			Phone: phone,
		})
	}
	otp = userOTP.OTP

	return otp, err
}

func (s *authService) SendOTPToken(cfg *config.RabbitMQ, emailPhone, channel, usedFor, otpToken, via string) error {
	// Define data mail
	dataMail := utils.DataEmailToken{ExpiredTime: model.TOKEN_EXPIRED_TIME}
	if channel == utils.ChannelApp {
		dataMail.Title = "OTP"
		dataMail.IsChannelApp = true
		dataMail.Description = "Please input the 6 digit code"
		dataMail.TokenURL = otpToken
	} else if channel == utils.ChannelWeb {
		dataMail.Title = "Link"
		dataMail.IsChannelApp = false
		dataMail.Description = "Please click the link"
		dataMail.TokenURL = fmt.Sprintf("https://sample.com/auth/%s?token=%s", usedFor, otpToken)
	}

	// Push email
	if via == model.OTP_VIA_EMAIL {
		return utils.PushQueueMail(cfg, utils.CMDQueueSendMail, utils.SetMailData([]string{emailPhone}, usedFor, dataMail))
	} else if via == model.OTP_VIA_PHONE {
		return errors.New("sorry, sms gateway not yet available")
	}

	return nil
}

func (s *authService) Registration(dbctx db.DBCtx, req requests.RegisterRequest, channel string, rabbitCfg *config.RabbitMQ) (model.User, string, error) {
	// Define data
	var user model.User
	var otpToken string

	// Compare password
	if req.Password != req.ConfirmPassword {
		return user, otpToken, fmt.Errorf("%s", "password is not match.")
	}

	// Phone validation
	_, _, validPhone := common.IsPhone(req.Phone)
	if !validPhone {
		return user, otpToken, fmt.Errorf(`phone %s is invalid`, req.Phone)
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return user, otpToken, err
	}

	// Set data
	code := common.CodeGenerator(model.USER_PREFIX, 5)
	user = model.User{
		Code:      code,
		Status:    false,
		Name:      req.Name,
		Email:     req.Email,
		Role:      model.ROLE_CUST,
		Phone:     req.Phone,
		Password:  string(passwordHash),
		CreatedBy: code,
	}

	// Set role slug
	roleID, err := s.roleR.GetIDBySlug(dbctx, user.Role)
	if err != nil {
		return user, otpToken, err
	}
	user.RoleID = int32(roleID)

	// Insert data
	userInserted, err := s.userR.Insert(dbctx, user)
	if err != nil {
		return user, otpToken, err
	}

	// Set / get otp user
	otpToken, err = s.GenerateOTPToken(dbctx, channel, userInserted.Email, userInserted.Phone)
	if err != nil {
		return user, otpToken, err
	}

	// Send otp to queue mail
	err = s.SendOTPToken(rabbitCfg, userInserted.Email, channel, utils.MAIL_FOR_USERACTIVATION, otpToken, model.OTP_VIA_EMAIL)
	if err != nil {
		return user, otpToken, err
	}

	return userInserted, otpToken, err
}

func (s *authService) AuthLogin(dbctx db.DBCtx, req requests.LoginRequest, channel string, rabbitCfg *config.RabbitMQ) (model.User, string, error) {
	// Define data
	var otpToken string

	// Get user
	user, err := s.userR.GetByEmail(dbctx, req.Email)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			err = fmt.Errorf(`invalid user credential`)
		}
		return user, otpToken, err
	}

	// Validate pass
	byteHash := []byte(user.Password)
	if err := bcrypt.CompareHashAndPassword(byteHash, []byte(req.Password)); err != nil {
		return user, otpToken, fmt.Errorf(`invalid user credential`)
	}

	// Adjustment login user role channel
	err = model.CheckValidChannelRole(channel, user.Role)
	if err != nil {
		return user, otpToken, err
	}

	// Adjustment user status
	if user.Status {
		// Generate token
		token, _, err := utils.GenerateJWT(user.ID, user.Code, user.Phone, user.Email, user.Role)
		if err != nil {
			return user, otpToken, err
		}
		user.RememberToken = sql.NullString{Valid: true, String: token}
	} else {
		// Set / get otp user
		otpToken, err = s.GenerateOTPToken(dbctx, channel, user.Email, user.Phone)
		if err != nil {
			return user, otpToken, err
		}

		// Send otp to queue mail
		err = s.SendOTPToken(rabbitCfg, user.Email, channel, utils.MAIL_FOR_USERACTIVATION, otpToken, model.OTP_VIA_EMAIL)
		if err != nil {
			return user, otpToken, err
		}
	}

	return user, otpToken, err
}

func (s *authService) SendOTPTokenByType(dbctx db.DBCtx, emailPhone, channel, sendType string, rabbitCfg *config.RabbitMQ) (string, bool, error) {
	// Define data
	var otpToken string

	// Check valid emailPhone
	via, validChannel := model.ViaValidMailPhoneChannel(channel, emailPhone)
	if !validChannel {
		return otpToken, false, errors.New("invalid channel email or phone")
	}

	// Valid type usage
	err := utils.CheckValidUsedFor(sendType)
	if err != nil {
		return otpToken, false, err
	}

	// Get user
	user, err := s.userR.GetByEmailOrPhone(dbctx, emailPhone)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			err = fmt.Errorf(`invalid user credential`)
		}
		return otpToken, user.Status, err
	}

	// Adjustment login user role channel
	err = model.CheckValidChannelRole(channel, user.Role)
	if err != nil {
		return otpToken, user.Status, err
	}

	// Adjustment sendType user not active
	if !user.Status {
		sendType = utils.MAIL_FOR_USERACTIVATION
	}

	// Set / get otp user
	otpToken, err = s.GenerateOTPToken(dbctx, channel, user.Email, user.Phone)
	if err != nil {
		return otpToken, user.Status, err
	}

	// Send otp to queue mail
	err = s.SendOTPToken(rabbitCfg, emailPhone, channel, sendType, otpToken, via)

	return otpToken, user.Status, err
}

func (s *authService) ChangePassword(dbctx db.DBCtx, req requests.ResetPasswordRequest, channel string) error {
	// Adjustment emailPhone channel
	emailPhone := req.EmailPhoneToken
	if channel == utils.ChannelWeb {
		// Decode token
		tokenDecoded, err := utils.DecodeTokenJWT(emailPhone)
		if err != nil {
			return err
		}
		emailPhone = tokenDecoded["email"].(string)
	}

	// Compare password
	if req.Password != req.ConfirmPassword {
		return fmt.Errorf("%s", "password is not match.")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	newPassword := string(passwordHash)

	// Update password
	err = s.userR.UpdatePasswordByEmailOrPhone(dbctx, newPassword, emailPhone)

	return err
}

func (s *authService) OTPTokenValidation(dbctx db.DBCtx, req requests.ValidateOTPTokenRequest, channel, sendType string) error {
	// Valid type usage
	err := utils.CheckValidUsedFor(sendType)
	if err != nil {
		return err
	}

	// Check valid emailPhone
	_, validChannel := model.ViaValidMailPhoneChannel(channel, req.EmailPhone)
	if !validChannel {
		return errors.New("invalid channel email or phone")
	}

	if channel == utils.ChannelWeb {
		// Decode token
		tokenDecoded, err := utils.DecodeTokenJWT(req.OTPToken)
		if err != nil {
			return err
		}

		// Check decode token with email
		if tokenDecoded["email"].(string) != req.EmailPhone {
			return errors.New("invalid token email")
		}
	} else if channel == utils.ChannelApp {
		// Get user otp
		userOTP, err := s.otpR.GetByEmailOrPhone(dbctx, req.EmailPhone, req.EmailPhone)
		if err != nil {
			return err
		}

		// Check expired time
		if time.Now().In(time.UTC).After(userOTP.ExpiredDate) {
			return fmt.Errorf("%s", "otp expired")
		}

		// Check contains equal otp
		if strings.Trim(userOTP.OTP, " ") != strings.Trim(req.OTPToken, " ") {
			return fmt.Errorf("%s", "invalid otp")
		}
	}

	// Update status active user if send type activation
	if strings.Trim(sendType, " ") == utils.MAIL_FOR_USERACTIVATION {
		return s.userR.UpdateStatusByEmailOrPhone(dbctx, true, req.EmailPhone)
	}

	return err
}
