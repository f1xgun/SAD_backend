package grades

import (
	"net/http"
	"sad/internal/handlers/grades"

	"github.com/go-chi/chi/v5"
)

func Routes(
	r *chi.Mux,
	handler grades.GradesHandler,
	authMiddleware func(http.Handler) http.Handler,
	allowedRolesMiddleware func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/grades/student/{student_id}", handler.GetAllStudentGrades)
		r.Get("/grades", handler.GetStudentGradesBySubjectAndGroup)

		r.Group(func(r chi.Router) {
			r.Use(allowedRolesMiddleware)
			r.Get("/grades/get_report_csv", handler.GetGradesInCsvReport)
			r.Post("/grades", handler.Create)
			r.Delete("/grades/{grade_id}", handler.Delete)
			r.Patch("/grades/{grade_id}", handler.Update)
		})
	})
}
