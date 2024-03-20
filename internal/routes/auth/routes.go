package routes

import (
	"sad/internal/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r *fiber.App, handler auth.AuthHandler) {
	r.Post("/login", handler.Login)
	r.Post("/register", handler.Register)
}
