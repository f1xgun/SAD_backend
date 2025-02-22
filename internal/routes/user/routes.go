package users

import (
	"net/http"
	users "sad/internal/handlers/user"

	"github.com/go-chi/chi/v5"
)

func Routes(
	r *chi.Mux, handler users.UserHandler,
	authMiddleware func(http.Handler) http.Handler,
	allowedRolesMiddleware func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/users/info", handler.GetUserInfoByToken)
		r.Get("/users/list", handler.GetUsers)
		r.Get("/users/info/{user_id}", handler.GetUserInfo)

		r.Group(func(r chi.Router) {
			r.Use(allowedRolesMiddleware)
			r.Patch("/users/{user_id}/edit", handler.EditUser)
			r.Delete("/users/{user_id}", handler.DeleteUser)
		})
	})
}
