package middleware

import (
	"fmt"
	"net/http"
	"sad/internal/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(config config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			tokenString = strings.TrimPrefix(authorization, "Bearer ")
		}

		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "You are not logged in"})
		}

		tokenByte, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.JwtSecret), nil
		})

		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": fmt.Sprintf("invalidate token: %v", err)})
		}

		if !tokenByte.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
		}

		return c.Next()
	}
}
