package subjects

import (
	"sad/internal/handlers/subjects"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App, handler subjects.Handler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	subjectApi := r.Group("/api/subjects").Use(authMiddleware)

	subjectApi.Get("/", handler.GetAll)                                                                                     // Получить список всех предметов
	subjectApi.Get("/get_by_teacher_id", handler.GetSubjectsByTeacherId)                                                    // Получить список предеметов, преподаваемых преподавателем
	subjectApi.Post("/", handler.Create).Use(allowedRolesMiddleware)                                                        // Создать новый предмет
	subjectApi.Get("/available_teachers", handler.GetAvailableTeachers)                                                     // Получить преподавателей
	subjectApi.Get("/get_new_available_for_teacher", handler.GetNewAvailableSubjectsForTeacher).Use(allowedRolesMiddleware) // Получить список предметов, которые можно назначить преподавателю
	subjectApi.Patch("/edit_teacher_subjects", handler.EditTeacherSubjects).Use(allowedRolesMiddleware)                     // Изменить список преподаваемых дисциплин у преподавателя

	subject := subjectApi.Group("/:subject_id")

	subject.Get("/details", handler.GetWithDetails)
	subject.Delete("/", handler.Delete).Use(allowedRolesMiddleware)    // Удалить предмет по ID
	subject.Patch("/edit", handler.Update).Use(allowedRolesMiddleware) // Обновить предмет по ID

	subject.Post("/groups", handler.AddSubjectToGroup).Use(allowedRolesMiddleware)        // Добавить предмет для группы
	subject.Delete("/groups", handler.DeleteSubjectFromGroup).Use(allowedRolesMiddleware) // Удалить предмет у группы
}
