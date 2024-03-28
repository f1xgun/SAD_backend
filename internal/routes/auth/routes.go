package auth

import (
	"sad/internal/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r *fiber.App, handler auth.AuthHandler) {
	r.Post("/api/login", handler.Login)
	r.Post("/api/register", handler.Register)
}
