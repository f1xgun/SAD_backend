package grades

import (
	"github.com/gofiber/fiber/v2"
	"sad/internal/handlers/grades"
)

func GradesRoutes(r *fiber.App, handler grades.GradesHandler, authMiddleware interface{}, allowedRoleMiddleware interface{}) {
	gradesApi := r.Group("/api/grades").Use(authMiddleware)

	userGradeApi := gradesApi.Group("/:student_id")
	userGradeApi.Get("/", handler.GetAllStudentGrades)

	gradesApi.Post("/", handler.Create).Use(allowedRoleMiddleware) // Создать новую оценку

	gradeApi := gradesApi.Group("/:grade_id")
	gradeApi.Delete("/", handler.Delete) // Удалить оценку по ID
	gradeApi.Patch("/", handler.Update)  // Обновить оценку по ID
}
