package groups

import (
	"net/http"
	"sad/internal/handlers/groups"

	"github.com/go-chi/chi/v5"
)

func Routes(
	r *chi.Mux,
	handler groups.Handler,
	authMiddleware func(http.Handler) http.Handler,
	allowedRolesMiddleware func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/groups", handler.GetAll)
		r.Get("/groups/teacher", handler.GetGroupsWithSubjectsByTeacher)
		r.Get("/groups/get_by_subject", handler.GetTeacherGroupsBySubject)

		r.Get("/groups/{group_id}", handler.Get)
		r.Get("/groups/{group_id}/details", handler.GetWithDetails)
		r.Group(func(r chi.Router) {
			r.Use(allowedRolesMiddleware)
			r.Post("/groups", handler.Create)
			r.Delete("/groups/{group_id}", handler.Delete)
			r.Patch("/groups/{group_id}", handler.Update)
			r.Get("/groups/{group_id}/available_new_users", handler.GetAvailableNewUsers)
			r.Post("/groups/{group_id}/users", handler.AddUserToGroup)
			r.Delete("/groups/{group_id}/users/{user_id}", handler.DeleteUserFromGroup)
		})
	})
}
