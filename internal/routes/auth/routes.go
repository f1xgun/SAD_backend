package auth

import (
	"sad/internal/handlers/auth"

	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux, handler auth.AuthHandler) {
	r.Post("/login", handler.Login)
	r.Post("/register", handler.Register)
}
