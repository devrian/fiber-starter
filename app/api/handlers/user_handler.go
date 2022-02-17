package handlers

import (
	"context"
	"fiber-starter/app/api"
	"fiber-starter/app/api/requests"
	"fiber-starter/app/api/responses"
	"fiber-starter/app/model"
	"fiber-starter/app/service"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"fiber-starter/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	app   *api.ApiApp
	userS service.UserService
	authS service.AuthService
}

func NewUserHandler(app *api.ApiApp, user service.UserService, auth service.AuthService) *UserHandler {
	return &UserHandler{app, user, auth}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.UserCreateRequest
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

	// Create user
	user, err := h.userS.CreateUser(dbctx, req, userData.Code, userData.Role)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set / get otp user
	otpToken, err := h.authS.GenerateOTPToken(dbctx, c.Get("X-Channel"), user.Email, user.Phone)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Send otp to queue mail
	err = h.authS.SendOTPToken(h.app.RabbitMQ, user.Email, c.Get("X-Channel"), utils.MAIL_FOR_USERACTIVATION, otpToken, model.OTP_VIA_EMAIL)
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
	var response responses.RegisterResponse
	response.Transform(user, otpToken)

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *UserHandler) GetList(c *fiber.Ctx) error {
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
	list, total, pg, err := h.userS.FindAllUser(dbctx, c)
	if err != nil {
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var listResp []responses.UserResponse
	for _, data := range list {
		var resp responses.UserResponse
		resp.Transform(data, "")
		listResp = append(listResp, resp)
	}

	pathResp := fmt.Sprintf(`%s?`, c.Route().Path)
	addResp := utils.WithPagination(listResp, total, common.OrderByOffsetLimitPaginateLink(pg, pathResp, nil))

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", addResp)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
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
	user, err := h.userS.FindUser(dbctx, c.Params("code"))
	if err != nil {
		return utils.APIResponse(c, db.ParseErr(err), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var response responses.UserResponse
	response.Transform(user, "")

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.UserUpdateRequest
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

	// Define code
	code := c.Params("code")
	if *req.IsProfile {
		code = userData.Code
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
	user, err := h.userS.UpdateUser(dbctx, req, userData.Code, code)
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
	var response responses.UserResponse
	response.Transform(user, "")

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
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
	err = h.userS.DeleteUser(dbctx, userData.Code, c.Params("code"))
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
