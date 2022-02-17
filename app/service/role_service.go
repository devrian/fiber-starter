package service

import (
	"database/sql"
	"errors"
	"fiber-starter/app/api/requests"
	"fiber-starter/app/model"
	"fiber-starter/app/repository"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
)

type RoleService interface {
	FindRole(dbctx db.DBCtx, code string) (model.Role, error)
	CreateRole(dbctx db.DBCtx, req requests.RoleCreateRequest, handlerBy string) (model.Role, error)
	UpdateRole(dbctx db.DBCtx, req requests.RoleUpdateRequest, handlerBy, code string) (model.Role, error)
	DeleteRole(dbctx db.DBCtx, handlerBy, code string) error
	FindAllRole(dbctx db.DBCtx, c *fiber.Ctx) ([]model.Role, int64, common.PaginateQueryOffset, error)
}

type roleService struct {
	roleR repository.RoleRepository
}

func NewRoleService(role repository.RoleRepository) *roleService {
	return &roleService{role}
}

func (s *roleService) FindRole(dbctx db.DBCtx, code string) (model.Role, error) {
	// Check code
	if len(code) <= 0 {
		return model.Role{}, errors.New("invalid code")
	}

	return s.roleR.GetByCode(dbctx, code)
}

func (s *roleService) CreateRole(dbctx db.DBCtx, req requests.RoleCreateRequest, handlerBy string) (model.Role, error) {
	// Insert role
	return s.roleR.Insert(dbctx, model.Role{
		Name:      req.Name,
		Status:    *req.Status,
		CreatedBy: handlerBy,
		Code:      common.CodeGenerator(model.ROLE_PREFIX, 5),
	})
}

func (s *roleService) UpdateRole(dbctx db.DBCtx, req requests.RoleUpdateRequest, handlerBy, code string) (model.Role, error) {
	// Check code
	if len(code) <= 0 {
		return model.Role{}, errors.New("invalid code")
	}

	// Define result
	var role model.Role

	// Check version
	version, err := s.roleR.GetVersionByCode(dbctx, code)
	if err != nil {
		return role, err
	}
	if version != int32(req.Version) {
		return role, errors.New("version is not match")
	}

	// Update role
	role = model.Role{
		Code:        code,
		Name:        req.Name,
		Status:      *req.Status,
		Version:     version + 1,
		Slug:        slug.Make(string(req.Name)),
		UpdatedDate: sql.NullTime{Valid: true, Time: time.Now().In(time.UTC)},
		UpdatedBy:   sql.NullString{Valid: true, String: handlerBy},
	}
	err = s.roleR.Update(dbctx, role)

	return role, err
}

func (s *roleService) DeleteRole(dbctx db.DBCtx, handlerBy, code string) error {
	// Check code
	if len(code) <= 0 {
		return errors.New("invalid code")
	}

	return s.roleR.Delete(dbctx, code, handlerBy)
}

func (s *roleService) FindAllRole(dbctx db.DBCtx, c *fiber.Ctx) ([]model.Role, int64, common.PaginateQueryOffset, error) {
	// Define variable
	var roleList []model.Role
	var total int64

	// Set filter
	var f model.RoleFilter
	search := c.Query("search")
	if len(search) > 0 {
		f.Search = search
	}
	status := c.Query("status")
	if len(status) > 0 {
		f.Status = status
	}

	// Define pagination
	pg, err := common.GetPaginateQueryOffset(c)
	if err != nil {
		return roleList, total, pg, err
	}
	f.Paging = pg

	// Get data
	roleList, err = s.roleR.GetAll(dbctx, f)
	if err != nil {
		return roleList, total, pg, err
	}

	// Get total data
	total, err = s.roleR.GetAllTotal(dbctx, f)

	return roleList, total, pg, err
}
