package groups

import (
	"sad/internal/handlers/groups"

	"github.com/gofiber/fiber/v2"
)

func GroupsRoutes(r *fiber.App, handler groups.GroupsHandler, authMiddlware interface{}, adminMiddleware interface{}) {
	groupApi := r.Group("/api/groups").Use(authMiddlware)

	groupApi.Get("/", handler.GetAll)                       // Получить список всех групп
	groupApi.Post("/", handler.Create).Use(adminMiddleware) // Создать новую группу

	group := groupApi.Group("/:group_id")

	group.Get("/", handler.Get)                            // Получить группу по ID
	group.Delete("/", handler.Delete).Use(adminMiddleware) // Удалить группу по ID
	group.Patch("/", handler.Update).Use(adminMiddleware)  // Обновить группу по ID

	group.Post("/users", handler.AddUserToGroup).Use(adminMiddleware)        // Добавить пользователя в группу
	group.Delete("/users", handler.DeleteUserFromGroup).Use(adminMiddleware) // Удалить пользователя из группы
}
