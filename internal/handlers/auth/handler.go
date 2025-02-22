package auth

import (
	"errors"
	"log/slog"
	"net/http"
	authModels "sad/internal/models/auth"
	errorsModels "sad/internal/models/errors"
	usersModels "sad/internal/models/users"
	"sad/internal/services"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService services.AuthService
	log         *slog.Logger
}

func NewAuthHandler(logger *slog.Logger, authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
		log:         logger,
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var user usersModels.UserCredentials

	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	token, err := h.authService.Login(r.Context(), user)

	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this login does not exist"
		case errors.Is(err, errorsModels.ErrInvalidCredentials):
			statusCode = http.StatusUnauthorized
			errMsg = "Invalid credentials"
		case errors.Is(err, errorsModels.ErrServer):
			statusCode = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = "An unexpected error occurred"
		}

		log.Error("login failed", slog.String("error", err.Error()))
		render.Status(r, statusCode)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("login successful", slog.String("login", user.Login))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"token": token})
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Register"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var user authModels.UserRegistrationRequest

	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("attempting to register user", slog.String("login", user.Login))

	err := h.authService.Register(r.Context(), user)
	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrUserExists):
			statusCode = http.StatusConflict
			errMsg = "User with this login already exists"
			log.Warn("registration failed: user already exists", slog.String("login", user.Login))
		case errors.Is(err, errorsModels.ErrServer):
			statusCode = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("registration failed: server error", slog.String("error", err.Error()))
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("registration failed: unexpected error", slog.String("error", err.Error()))
		}

		render.Status(r, statusCode)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("registration successful", slog.String("login", user.Login))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "registration successful"})
}
