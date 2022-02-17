package routes

import (
	"fiber-starter/app/api/handlers"
	"fiber-starter/app/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type PublicHandlers struct {
	Auth *handlers.AuthHandler
}

// PublicRoutes func for describe group of public routes.
func PublicRoutes(r fiber.Router, h PublicHandlers) {
	// Route Auth
	auth := r.Group("/auth")
	auth.Post("/register", middleware.ChannelAppOnly(), h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/send-otptoken/:type?", h.Auth.SendOTPToken)
	auth.Post("/reset-password", h.Auth.ResetPassword)
	auth.Post("/validate-otptoken/:type?", h.Auth.ValidateOTPToken)
}
