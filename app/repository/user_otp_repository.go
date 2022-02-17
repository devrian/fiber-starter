package repository

import (
	"fiber-starter/app/model"
	"fiber-starter/db"
	"time"
)

type UserOTPRepository interface {
	Insert(dbctx db.DBCtx, u model.UserOTP) (model.UserOTP, error)
	UpdateOTP(dbctx db.DBCtx, u model.UserOTP) (model.UserOTP, error)
	GetByEmailOrPhone(dbctx db.DBCtx, email, phone string) (model.UserOTP, error)
}

type userOTPRepository struct {
}

func NewUserOTPRepository() *userOTPRepository {
	return &userOTPRepository{}
}

func (r *userOTPRepository) Insert(dbctx db.DBCtx, u model.UserOTP) (model.UserOTP, error) {
	var ID int64
	var timeStamp = time.Now().In(time.UTC)

	otp, err := u.GenerateOTP()
	if err != nil {
		return u, err
	}
	u.OTP = otp
	u.ExpiredDate = u.GetExpiredDate()

	paramQ := []interface{}{u.Email, u.Phone, u.OTP, u.ExpiredDate, timeStamp}
	q := `insert into user_otps (email, phone, otp, expired_date, created_date) values ($1, $2, $3, $4, $5) returning id`
	err = dbctx.TX.QueryRow(dbctx.Ctx, q, paramQ...).Scan(&ID)
	if err != nil {
		return u, err
	}
	u.ID = ID

	return u, err
}

func (r *userOTPRepository) UpdateOTP(dbctx db.DBCtx, u model.UserOTP) (model.UserOTP, error) {
	var ID int64
	var timeStamp = time.Now().In(time.UTC)

	otp, err := u.GenerateOTP()
	if err != nil {
		return u, err
	}
	u.OTP = otp
	u.ExpiredDate = u.GetExpiredDate()

	paramQ := []interface{}{u.OTP, u.ExpiredDate, timeStamp, u.Email, u.Phone}
	q := `update user_otps set otp = $1, expired_date = $2, updated_date = $3 where (email = $4 or phone = $5) returning id`
	err = dbctx.TX.QueryRow(dbctx.Ctx, q, paramQ...).Scan(&ID)
	if err != nil {
		return u, err
	}
	u.ID = ID

	return u, err
}

func (r userOTPRepository) GetByEmailOrPhone(dbctx db.DBCtx, email, phone string) (model.UserOTP, error) {
	var u model.UserOTP

	q := `select id, email, phone, otp, expired_date from user_otps where (email = $1 or phone = $2)`
	err := dbctx.DB.QueryRow(dbctx.Ctx, q, email, phone).Scan(&u.ID, &u.Email, &u.Phone, &u.OTP, &u.ExpiredDate)

	return u, err
}
