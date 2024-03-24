package user

import (
	"sad/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r *fiber.App, handler user.UserHandler, middlewares ...interface{}) {
	api := r.Group("/user").Use(middlewares...)
	api.Patch("/role", handler.EditRole)
}
