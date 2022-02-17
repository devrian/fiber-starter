package middleware

import (
	"fiber-starter/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(app *fiber.App) {
	allowHeaders := []string{
		"Origin",
		"Accept",
		"Authorization",
		"Content-Type",
		"X-CSRF-Token",
		"X-SIGNATURE",
		"X-TIMESTAMPT",
		"X-CHANNEL",
		"X-PLAYER",
		"Access-Control-Allow-Headers",
		"X-Requested-With",
		"application/json",
		"Cache-Control",
		"Token",
		"X-Token",
	}

	app.Use(
		// Add CORS to each route.
		cors.New(cors.Config{
			AllowHeaders: strings.Join(allowHeaders, ", "),
		}),
		// Add simple logger.
		logger.New(),
	)
}

func ChannelAppOnly() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Get("X-Channel") == utils.ChannelApp
		},
	})
}
