package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type TokenMetaData struct {
	ID      int64
	Code    string
	Phone   string
	Email   string
	Role    string
	Expires int64
}

func GenerateJWT(id int64, code, phone, email, role string) (string, int64, error) {
	// Set expired time
	expirationTime := time.Now().Add(2160 * time.Hour).Unix()

	// Set secret key from .env file.
	appKey := os.Getenv("APP_KEY")

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["id"] = id
	claims["code"] = code
	claims["phone"] = phone
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = expirationTime

	// Create a new JWT access token with claims.
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(appKey))

	return token, expirationTime, err
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetaData, error) {
	var err error
	var tokenMetaData *TokenMetaData

	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		tokenMetaData = &TokenMetaData{
			ID:      int64(claims["id"].(float64)),
			Code:    fmt.Sprintf("%s", claims["code"]),
			Phone:   fmt.Sprintf("%s", claims["phone"]),
			Email:   fmt.Sprintf("%s", claims["email"]),
			Role:    fmt.Sprintf("%s", claims["role"]),
			Expires: int64(claims["exp"].(float64)),
		}
	}

	return tokenMetaData, err
}

func DecodeTokenJWT(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(tokenTemp *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("APP_KEY")), nil
	})
	if err != nil {
		return claims, err
	}

	return claims, nil
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("APP_KEY")), nil
}
