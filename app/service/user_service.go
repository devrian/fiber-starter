package service

import (
	"database/sql"
	"errors"
	"fiber-starter/app/api/requests"
	"fiber-starter/app/model"
	"fiber-starter/app/repository"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindUser(dbctx db.DBCtx, code string) (model.User, error)
	FindAllUser(dbctx db.DBCtx, c *fiber.Ctx) ([]model.User, int64, common.PaginateQueryOffset, error)
	CreateUser(dbctx db.DBCtx, req requests.UserCreateRequest, handlerBy, roleBy string) (model.User, error)
	UpdateUser(dbctx db.DBCtx, req requests.UserUpdateRequest, handlerBy, code string) (model.User, error)
	DeleteUser(dbctx db.DBCtx, handlerBy, code string) error
}

type userService struct {
	userR repository.UserRepository
	roleR repository.RoleRepository
	otpR  repository.UserOTPRepository
}

func NewUserService(user repository.UserRepository, role repository.RoleRepository, otp repository.UserOTPRepository) *userService {
	return &userService{user, role, otp}
}

func (s *userService) FindUser(dbctx db.DBCtx, code string) (model.User, error) {
	// Check code
	if len(code) <= 0 {
		return model.User{}, errors.New("invalid code")
	}

	return s.userR.GetByCode(dbctx, code)
}

func (s *userService) FindAllUser(dbctx db.DBCtx, c *fiber.Ctx) ([]model.User, int64, common.PaginateQueryOffset, error) {
	// Define variable
	var userList []model.User
	var total int64

	// Set filter
	var f model.UserFilter
	search := c.Query("search")
	if len(search) > 0 {
		f.Search = search
	}
	status := c.Query("status")
	if len(status) > 0 {
		f.Status = status
	}
	roles := c.Query("roles")
	if len(roles) > 0 {
		listRoles := strings.Split(roles, ",")
		if len(listRoles) > 0 {
			for _, role := range listRoles {
				if err := model.CheckValidRoleSlug(role); err != nil {
					return userList, total, common.PaginateQueryOffset{}, err
				}
			}
			f.Roles = listRoles
		}
	}

	// Define pagination
	pg, err := common.GetPaginateQueryOffset(c)
	if err != nil {
		return userList, total, pg, err
	}
	f.Paging = pg

	// Get data
	userList, err = s.userR.GetAll(dbctx, f)
	if err != nil {
		return userList, total, pg, err
	}

	// Get total data
	total, err = s.userR.GetAllTotal(dbctx, f)

	return userList, total, pg, err
}

func (s *userService) CreateUser(dbctx db.DBCtx, req requests.UserCreateRequest, handlerBy, roleBy string) (model.User, error) {
	// Define data
	var user model.User

	// Check role by
	if !common.CheckStringContains(roleBy, []string{model.ROLE_ADMIN}) {
		return user, errors.New("invalid role access")
	}

	// Phone validation
	_, _, validPhone := common.IsPhone(req.Phone)
	if !validPhone {
		return user, fmt.Errorf(`phone %s is invalid`, req.Phone)
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(model.DEFAULT_PASS), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	// Set data
	code := common.CodeGenerator(model.USER_PREFIX, 5)
	user = model.User{
		Code:      code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(passwordHash),
		Address:   sql.NullString{Valid: true, String: req.Address},
		Img:       sql.NullString{Valid: true, String: req.Img},
		CreatedBy: handlerBy,
	}

	if len(req.RoleCode) <= 0 {
		return user, errors.New("role code is required")
	}
	role, err := s.roleR.GetByCode(dbctx, req.RoleCode)
	if err != nil {
		return user, err
	}
	user.RoleID = int32(role.ID)
	user.Role = role.Slug

	// Insert data
	newUser, err := s.userR.Insert(dbctx, user)
	if err != nil {
		return user, err
	}

	return newUser, err
}

func (s *userService) UpdateUser(dbctx db.DBCtx, req requests.UserUpdateRequest, handlerBy, code string) (model.User, error) {
	// Check code
	if len(code) <= 0 {
		return model.User{}, errors.New("invalid code")
	}

	// Define result
	var user model.User

	// Check version
	version, err := s.userR.GetVersionByCode(dbctx, code)
	if err != nil {
		return user, err
	}
	if version != int32(req.Version) {
		return user, errors.New("version is not match")
	}

	// Update role
	user = model.User{
		Code:        code,
		Name:        req.Name,
		Phone:       req.Phone,
		Address:     sql.NullString{Valid: true, String: req.Address},
		Img:         sql.NullString{Valid: true, String: req.Img},
		Status:      *req.Status,
		Version:     version + 1,
		UpdatedDate: sql.NullTime{Valid: true, Time: time.Now().In(time.UTC)},
		UpdatedBy:   sql.NullString{Valid: true, String: handlerBy},
	}
	err = s.userR.Update(dbctx, user)

	return user, err
}

func (s *userService) DeleteUser(dbctx db.DBCtx, handlerBy, code string) error {
	// Check code
	if len(code) <= 0 {
		return errors.New("invalid code")
	}

	return s.userR.Delete(dbctx, code, handlerBy)
}
