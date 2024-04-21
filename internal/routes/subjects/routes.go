package subjects

import (
	"sad/internal/handlers/subjects"

	"github.com/gofiber/fiber/v2"
)

func SubjectsRoutes(r *fiber.App, handler subjects.SubjectsHandler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	subjectApi := r.Group("/api/subjects").Use(authMiddleware)

	subjectApi.Get("/", handler.GetAll)                              // Получить список всех предметов
	subjectApi.Post("/", handler.Create).Use(allowedRolesMiddleware) // Создать новый предмет

	subject := subjectApi.Group("/:subject_id")

	subject.Delete("/", handler.Delete).Use(allowedRolesMiddleware) // Удалить предмет по ID
	subject.Patch("/", handler.Update).Use(allowedRolesMiddleware)  // Обновить предмет по ID

	subject.Post("/groups", handler.AddSubjectToGroup).Use(allowedRolesMiddleware)        // Добавить предмет для группы
	subject.Delete("/groups", handler.DeleteSubjectFromGroup).Use(allowedRolesMiddleware) // Удалить предмет у группы
}
