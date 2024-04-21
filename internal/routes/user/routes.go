package users

import (
	users "sad/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r *fiber.App, handler users.UserHandler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	userApi := r.Group("/api/users").Use(authMiddleware)

	userApi.Get("/info", handler.GetUserInfo)
	userApi.Patch("/:user_id/role", handler.EditRole).Use(allowedRolesMiddleware)
}
