package users

import (
	users "sad/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r *fiber.App, handler users.UserHandler, authMiddlware interface{}, adminMiddleware interface{}) {
	userApi := r.Group("/api/users")

	userApi.Patch("/:user_id/role", handler.EditRole).Use(authMiddlware, adminMiddleware)
}
