package handlers

import (
	"context"
	"fiber-starter/app/api"
	"fiber-starter/app/api/requests"
	"fiber-starter/app/api/responses"
	"fiber-starter/app/service"
	"fiber-starter/db"
	"fiber-starter/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	app   *api.ApiApp
	authS service.AuthService
}

func NewAuthHandler(app *api.ApiApp, auth service.AuthService) *AuthHandler {
	return &AuthHandler{app, auth}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.RegisterRequest
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

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Registration
	user, otpToken, err := h.authS.Registration(dbctx, req, c.Get("X-Channel"), h.app.RabbitMQ)
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

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.LoginRequest
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

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Login
	user, otpToken, err := h.authS.AuthLogin(dbctx, req, c.Get("X-Channel"), h.app.RabbitMQ)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var response responses.LoginResponse
	response.Transform(user, otpToken)

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *AuthHandler) SendOTPToken(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.ResetRequest
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

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Forgot assword
	otpToken, status, err := h.authS.SendOTPTokenByType(dbctx, req.EmailPhone, c.Get("X-Channel"), c.Params("type"), h.app.RabbitMQ)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Set response
	var response responses.SendOTPTokenResponse
	response.Transform(otpToken, status, req.EmailPhone)

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", response)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.ResetPasswordRequest
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

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Change password
	err = h.authS.ChangePassword(dbctx, req, c.Get("X-Channel"))
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	return utils.APIResponse(c, "success", fiber.StatusOK, "success", nil)
}

func (h *AuthHandler) ValidateOTPToken(c *fiber.Ctx) error {
	// Define request with validation
	var req requests.ValidateOTPTokenRequest
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

	// Set Tx transaction
	tx, err := h.app.DB.Begin(ctx)
	if err != nil {
		return utils.APIResponse(c, err.Error(), fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error(), nil)
	}

	// Set db context
	var dbctx db.DBCtx
	dbctx.Set(ctx, conn, tx)

	// Validate
	err = h.authS.OTPTokenValidation(dbctx, req, c.Get("X-Channel"), c.Params("type"))
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
