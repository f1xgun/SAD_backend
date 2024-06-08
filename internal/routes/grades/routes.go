package grades

import (
	"github.com/gofiber/fiber/v2"
	"sad/internal/handlers/grades"
)

func Routes(r *fiber.App, handler grades.Handler, authMiddleware interface{}, allowedRoleMiddleware interface{}) {
	gradesApi := r.Group("/api/grades").Use(authMiddleware)
	gradesApi.Get("/", handler.GetStudentGradesBySubjectAndGroup)

	userGradeApi := gradesApi.Group("/:student_id")
	userGradeApi.Get("/", handler.GetAllStudentGrades)
	//userGradeApi.Get("/final", handler.GetAllStudentFinalGrades)

	gradesApi.Post("/", handler.Create).Use(allowedRoleMiddleware) // Создать новую оценку

	gradeApi := gradesApi.Group("/:grade_id")
	gradeApi.Delete("/", handler.Delete) // Удалить оценку по ID
	gradeApi.Patch("/", handler.Update)  // Обновить оценку по ID
}
