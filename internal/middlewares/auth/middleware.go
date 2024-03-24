package auth

import (
	"fmt"
	"log"
	"net/http"
	"sad/internal/config"
	authModels "sad/internal/models/auth"
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
			log.Println("Authorization header is missing or does not contain Bearer token")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "You are not logged in"})
		}

		token, err := jwt.ParseWithClaims(tokenString, &authModels.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				errMsg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
				log.Println(errMsg)
				return nil, fmt.Errorf(errMsg)
			}

			return []byte(config.JwtSecret), nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
		}

		if claims, ok := token.Claims.(*authModels.Claims); ok && token.Valid {
			c.Locals("userID", claims.Subject)
			log.Printf("Token is valid for userID: %s", claims.Subject)
			return c.Next()
		} else {
			log.Println("Token is not valid or claims are not of expected type")
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"error": "unauthorized"})
		}
	}
}
