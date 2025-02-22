package subjects

import (
	"net/http"
	"sad/internal/handlers/subjects"

	"github.com/go-chi/chi/v5"
)

func Routes(
	r *chi.Mux,
	handler subjects.Handler,
	authMiddleware func(http.Handler) http.Handler,
	allowedRolesMiddleware func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/subjects", handler.GetAll)
		r.Get("/subjects/get_by_teacher_id", handler.GetSubjectsByTeacherId)

		r.Get("/subjects/{subject_id}/details", handler.GetWithDetails)
		r.Group(func(r chi.Router) {
			r.Use(allowedRolesMiddleware)
			r.Post("/subjects", handler.Create)
			r.Get("/subjects/available_teachers", handler.GetAvailableTeachers)
			r.Get("/subjects/get_new_available_for_teacher", handler.GetNewAvailableSubjectsForTeacher)
			r.Patch("/subjects/edit_teacher_subjects", handler.EditTeacherSubjects)

			r.Delete("/subjects/{subject_id}", handler.Delete)
			r.Patch("/subjects/{subject_id}/edit", handler.Update)
			r.Post("/subjects/{subject_id}/groups", handler.AddSubjectToGroup)
			r.Delete("/subjects/{subject_id}/groups", handler.DeleteSubjectFromGroup)
		})
	})
}
