package routes

import (
	"fiber-starter/app/api/handlers"
	"fiber-starter/app/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type PrivateHandlers struct {
	User *handlers.UserHandler
	Role *handlers.RoleHandler
}

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(r fiber.Router, h PrivateHandlers) {
	// Route Role
	role := r.Group("/role", middleware.JWTRoleAdmin())
	role.Post("/", h.Role.Create)
	role.Get("/", h.Role.GetList)
	role.Get("/:code?", h.Role.Get)
	role.Put("/:code?", h.Role.Update)
	role.Delete("/:code?", h.Role.Delete)

	// Route User
	user := r.Group("/user", middleware.JWTProtected())
	user.Post("/", h.User.Create)
	user.Get("/", h.User.GetList)
	user.Post("profile", h.User.Update)
	user.Get("/:code?", h.User.Get)
	user.Put("/:code?", h.User.Update)
	user.Delete("/:code?", h.User.Delete)
}
