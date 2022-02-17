package routes

import (
	"fiber-starter/app/api"
	"fiber-starter/app/api/handlers"
	"fiber-starter/app/api/middleware"
	"fiber-starter/app/repository"
	"fiber-starter/app/service"
	"fmt"
)

func Configure(app *api.ApiApp) {
	// Define Repositories
	userR := repository.NewUserRepository()
	roleR := repository.NewRoleRepository()
	otpR := repository.NewUserOTPRepository()

	// Define Services
	authS := service.NewAuthService(userR, roleR, otpR)
	roleS := service.NewRoleService(roleR)
	userS := service.NewUserService(userR, roleR, otpR)

	// Define Handlers
	authH := handlers.NewAuthHandler(app, authS)
	roleH := handlers.NewRoleHandler(app, roleS)
	userH := handlers.NewUserHandler(app, userS, authS)

	// Define Main Route API
	api := app.Fiber.Group(fmt.Sprintf("/api/%s", app.Config.App.Version))

	// Middlewares.
	middleware.FiberMiddleware(app.Fiber) // Register Fiber's middleware for app.

	// Routes
	PublicRoutes(api, PublicHandlers{authH})
	PrivateRoutes(api, PrivateHandlers{userH, roleH})
}
