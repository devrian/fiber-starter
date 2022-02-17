package middleware

import (
	"fiber-starter/app/model"
	"fiber-starter/pkg/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	jwtMiddleware "github.com/gofiber/jwt/v2"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("APP_KEY")),
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

// Role admin protection
func JWTRoleAdmin() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Next: func(c *fiber.Ctx) bool {
			tokenMetaData, err := utils.ExtractTokenMetadata(c)
			if err != nil {
				return false
			}
			return tokenMetaData.Role == model.ROLE_ADMIN
		},
	})
}

// Role customer protection
func JWTRoleCust() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Next: func(c *fiber.Ctx) bool {
			tokenMetaData, err := utils.ExtractTokenMetadata(c)
			if err != nil {
				return false
			}
			return tokenMetaData.Role == model.ROLE_CUST
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return utils.APIResponse(c, err.Error(), fiber.StatusBadRequest, fiber.ErrBadRequest.Error(), nil)
	}

	// Return status 401 and failed authentication error.
	return utils.APIResponse(c, err.Error(), fiber.StatusUnauthorized, fiber.ErrUnauthorized.Error(), nil)
}
