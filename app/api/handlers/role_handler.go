package handlers

import (
	"context"
	"fiber-starter/app/api"
	"fiber-starter/app/api/requests"
	"fiber-starter/app/api/responses"
	"fiber-starter/app/service"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"fiber-starter/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	app   *api.ApiApp
	roleS service.RoleService
}

func NewRoleHandler(app *api.ApiApp, role service.RoleService) *RoleHandler {
	return &RoleHandler{app, role}
}

func (h *RoleHandler) Create(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.RoleCreateRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	if err := h.app.Validator.Driver.Struct(req); err != nil {
		return utils.APIResponseErrorByValidationError(c, err)
	}

	// Set context
	ctx := context.Background()
	conn, err := h.app.DB.Acquire(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	defer conn.Release()

	// Get user code (handler by)
	userData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Create role
	role, err := h.roleS.CreateRole(dbctx, req, userData.Code)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var response responses.RoleResponse
	response.Transform(role)

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *RoleHandler) GetList(c *fiber.Ctx) error {
	// Set context
	ctx := context.Background()
	conn, err := h.app.DB.Acquire(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	defer conn.Release()

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, nil)

	// Get data
	list, total, pg, err := h.roleS.FindAllRole(dbctx, c)
	if err != nil {
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var listResp []responses.RoleResponse
	for _, data := range list {
		var resp responses.RoleResponse
		resp.Transform(data)
		listResp = append(listResp, resp)
	}

	pathResp := fmt.Sprintf(`%s?`, c.Route().Path)
	addResp := utils.WithPagination(listResp, total, common.OrderByOffsetLimitPaginateLink(pg, pathResp, nil))

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", addResp)
}

func (h *RoleHandler) Get(c *fiber.Ctx) error {
	// Set context
	ctx := context.Background()
	conn, err := h.app.DB.Acquire(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	defer conn.Release()

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, nil)

	// Find data
	role, err := h.roleS.FindRole(dbctx, c.Params("code"))
	if err != nil {
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var response responses.RoleResponse
	response.Transform(role)

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *RoleHandler) Update(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.RoleUpdateRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	if err := h.app.Validator.Driver.Struct(req); err != nil {
		return utils.APIResponseErrorByValidationError(c, err)
	}

	// Set context
	ctx := context.Background()
	conn, err := h.app.DB.Acquire(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	defer conn.Release()

	// Get user code (handler by)
	userData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Update role
	role, err := h.roleS.UpdateRole(dbctx, req, userData.Code, c.Params("code"))
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var response responses.RoleResponse
	response.Transform(role)

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *RoleHandler) Delete(c *fiber.Ctx) error {
	// Set context
	ctx := context.Background()
	conn, err := h.app.DB.Acquire(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}
	defer conn.Release()

	// Get user code (handler by)
	userData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Delete role
	err = h.roleS.DeleteRole(dbctx, userData.Code, c.Params("code"))
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", nil)
}
