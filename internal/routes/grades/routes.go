package grades

import (
	"github.com/gofiber/fiber/v2"
	"sad/internal/handlers/grades"
)

func Routes(r *fiber.App, handler grades.Handler, authMiddleware interface{}, allowedRoleMiddleware interface{}) {
	gradesBaseApi := r.Group("/api/grades").Use(authMiddleware)

	gradesApi := gradesBaseApi.Group("/")
	gradesApi.Get("/", handler.GetStudentGradesBySubjectAndGroup)

	userGradeApi := gradesBaseApi.Group("/student/:student_id")
	userGradeApi.Get("/", handler.GetAllStudentGrades)

	gradesAllowedRoleApi := gradesBaseApi.Group("/").Use(allowedRoleMiddleware)
	gradesAllowedRoleApi.Post("/", handler.Create) // Создать новую оценку
	gradesAllowedRoleApi.Get("/get_report_csv", handler.GetGradesInCsvReport)

	gradeApi := gradesAllowedRoleApi.Group("/:grade_id")
	gradeApi.Delete("/", handler.Delete) // Удалить оценку по ID
	gradeApi.Patch("/", handler.Update)  // Обновить оценку по ID
}
