package subjects

import (
	"sad/internal/handlers/subjects"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App, handler subjects.Handler, authMiddleware interface{}, allowedRolesMiddleware interface{}) {
	subjectBaseApi := r.Group("/api/subjects").Use(authMiddleware)

	subjectBaseApi.Get("/", handler.GetAll)                                  // Получить список всех предметов
	subjectBaseApi.Get("/get_by_teacher_id", handler.GetSubjectsByTeacherId) // Получить список предеметов, преподаваемых преподавателем

	subjectApi := subjectBaseApi.Group("/:subject_id")
	subjectApi.Get("/details", handler.GetWithDetails)

	subjectAllowedRolesApi := subjectBaseApi.Group("/").Use(allowedRolesMiddleware)
	subjectAllowedRolesApi.Post("/", handler.Create)                                                        // Создать новый предмет
	subjectAllowedRolesApi.Get("/available_teachers", handler.GetAvailableTeachers)                         // Получить преподавателей
	subjectAllowedRolesApi.Get("/get_new_available_for_teacher", handler.GetNewAvailableSubjectsForTeacher) // Получить список предметов, которые можно назначить преподавателю
	subjectAllowedRolesApi.Patch("/edit_teacher_subjects", handler.EditTeacherSubjects)                     // Изменить список преподаваемых дисциплин у преподавателя

	subjectApi.Delete("/", handler.Delete)    // Удалить предмет по ID
	subjectApi.Patch("/edit", handler.Update) // Обновить предмет по ID

	subjectApi.Post("/groups", handler.AddSubjectToGroup)        // Добавить предмет для группы
	subjectApi.Delete("/groups", handler.DeleteSubjectFromGroup) // Удалить предмет у группы
}
