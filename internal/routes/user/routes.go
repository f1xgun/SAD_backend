package users

import (
	users "sad/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App, handler users.UserHandler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	usersApi := r.Group("/api/users").Use(authMiddleware)
	usersApi.Get("/list", handler.GetUsers)
	usersApi.Get("/info", handler.GetUserInfoByToken)

	userApi := usersApi.Group("/:user_id")
	userApi.Get("/info", handler.GetUserInfo)
	userApi.Patch("/edit", handler.EditUser).Use(allowedRolesMiddleware)
	userApi.Delete("/", handler.DeleteUser).Use(allowedRolesMiddleware)
}
