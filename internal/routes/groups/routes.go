package groups

import (
	"sad/internal/handlers/groups"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App, handler groups.Handler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	groupBaseApi := r.Group("/api/groups").Use(authMiddleware)

	groupBaseApi.Get("/", handler.GetAll) // Получить список всех групп
	groupBaseApi.Get("/teacher", handler.GetGroupsWithSubjectsByTeacher)
	groupBaseApi.Get("/get_by_subject", handler.GetTeacherGroupsBySubject)

	groupApi := groupBaseApi.Group("/:group_id")

	groupApi.Get("/", handler.Get)                   // Получить группу по ID
	groupApi.Get("/details", handler.GetWithDetails) // Получить группу по ID с деталями

	groupAllowedRolesApi := groupApi.Group("/")
	groupAllowedRolesApi.Use(allowedRolesMiddleware)
	groupAllowedRolesApi.Post("/", handler.Create)   // Создать новую группу
	groupAllowedRolesApi.Delete("/", handler.Delete) // Удалить группу по ID
	groupAllowedRolesApi.Patch("/", handler.Update)  // Обновить группу по ID
	groupAllowedRolesApi.Get("/available_new_users", handler.GetAvailableNewUsers)

	groupAllowedRolesApi.Post("/users/", handler.AddUserToGroup)                // Добавить пользователя в группу
	groupAllowedRolesApi.Delete("/users/:user_id", handler.DeleteUserFromGroup) // Удалить пользователя из группы
}
