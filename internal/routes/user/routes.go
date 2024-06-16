package users

import (
	users "sad/internal/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App, handler users.UserHandler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	usersBaseApi := r.Group("/api/users").Use(authMiddleware)
	usersBaseApi.Get("/list", handler.GetUsers)
	usersBaseApi.Get("/info", handler.GetUserInfoByToken)

	userApi := usersBaseApi.Group("/:user_id")
	userApi.Get("/info", handler.GetUserInfo)

	userAllowedRolesApi := userApi.Use(allowedRolesMiddleware)
	userAllowedRolesApi.Patch("/edit", handler.EditUser)
	userAllowedRolesApi.Delete("/", handler.DeleteUser)
}
