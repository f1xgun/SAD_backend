package groups

import (
	"sad/internal/handlers/groups"

	"github.com/gofiber/fiber/v2"
)

func GroupsRoutes(r *fiber.App, handler groups.GroupsHandler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	groupApi := r.Group("/api/groups").Use(authMiddleware)

	groupApi.Get("/", handler.GetAll)                              // Получить список всех групп
	groupApi.Post("/", handler.Create).Use(allowedRolesMiddleware) // Создать новую группу

	group := groupApi.Group("/:group_id")

	group.Get("/", handler.Get)                   // Получить группу по ID
	group.Get("/details", handler.GetWithDetails) // Получить группу по ID с деталями

	groupWithAllowedRole := group.Use(allowedRolesMiddleware)
	groupWithAllowedRole.Delete("/", handler.Delete) // Удалить группу по ID
	groupWithAllowedRole.Patch("/", handler.Update)  // Обновить группу по ID
	groupWithAllowedRole.Get("/available_new_users", handler.GetAvailableNewUsers)

	groupWithAllowedRole.Post("/users/", handler.AddUserToGroup)                // Добавить пользователя в группу
	groupWithAllowedRole.Delete("/users/:user_id", handler.DeleteUserFromGroup) // Удалить пользователя из группы
}
