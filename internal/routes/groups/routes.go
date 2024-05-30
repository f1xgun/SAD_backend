package groups

import (
	"sad/internal/handlers/groups"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App, handler groups.Handler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	groupApi := r.Group("/api/groups").Use(authMiddleware)

	groupApi.Get("/", handler.GetAll) // Получить список всех групп
	groupApi.Get("/teacher", handler.GetGroupsWithSubjectsByTeacher)
	groupApi.Get("/get_by_subject", handler.GetTeacherGroupsBySubject)

	group := groupApi.Group("/:group_id")

	group.Get("/", handler.Get)                                    // Получить группу по ID
	group.Get("/details", handler.GetWithDetails)                  // Получить группу по ID с деталями
	groupApi.Post("/", handler.Create).Use(allowedRolesMiddleware) // Создать новую группу

	group.Delete("/", handler.Delete) // Удалить группу по ID
	group.Patch("/", handler.Update)  // Обновить группу по ID
	group.Get("/available_new_users", handler.GetAvailableNewUsers)

	group.Post("/users/", handler.AddUserToGroup)                // Добавить пользователя в группу
	group.Delete("/users/:user_id", handler.DeleteUserFromGroup) // Удалить пользователя из группы
}
