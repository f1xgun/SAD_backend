package subjects

import (
	"sad/internal/handlers/subjects"

	"github.com/gofiber/fiber/v2"
)

func SubjectsRoutes(r *fiber.App, handler subjects.SubjectsHandler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	subjectApi := r.Group("/api/subjects").Use(authMiddleware)

	subjectApi.Get("/", handler.GetAll)                                 // Получить список всех предметов
	subjectApi.Post("/", handler.Create).Use(allowedRolesMiddleware)    // Создать новый предмет
	subjectApi.Get("/available_teachers", handler.GetAvailableTeachers) // Получить преподавателей

	subject := subjectApi.Group("/:subject_id")

	subject.Get("/details", handler.GetWithDetails)
	subject.Delete("/", handler.Delete).Use(allowedRolesMiddleware)    // Удалить предмет по ID
	subject.Patch("/edit", handler.Update).Use(allowedRolesMiddleware) // Обновить предмет по ID

	subject.Post("/groups", handler.AddSubjectToGroup).Use(allowedRolesMiddleware)        // Добавить предмет для группы
	subject.Delete("/groups", handler.DeleteSubjectFromGroup).Use(allowedRolesMiddleware) // Удалить предмет у группы
}
