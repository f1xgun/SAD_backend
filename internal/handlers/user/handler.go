package users

import (
	"errors"
	"log/slog"
	"net/http"
	errorsModels "sad/internal/models/errors"
	"sad/internal/models/users"
	"sad/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserHandler interface {
	EditUser(w http.ResponseWriter, r *http.Request)
	GetUserInfo(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetUserInfoByToken(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService services.UserService
	log         *slog.Logger
}

func NewUserHandler(logger *slog.Logger, userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
		log:         logger,
	}
}

func (h *userHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.EditUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		log.Error("failed to get userID from context")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	var user users.UserInfo
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("attempting to edit user info", slog.String("user_id", userID))

	err := h.userService.EditUser(r.Context(), userID, user.Role, user.Name)
	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrNoPermission):
			statusCode = http.StatusForbidden
			errMsg = "No permission to change user info"
			log.Warn("no permission to change user info", slog.String("user_id", userID))
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
			log.Warn("user not found", slog.String("user_id", userID))
		case errors.Is(err, errorsModels.ErrChangeOwnRole):
			statusCode = http.StatusForbidden
			errMsg = "Cannot change own info"
			log.Warn("attempt to change own info", slog.String("user_id", userID))
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("failed to edit user info", slog.String("error", err.Error()))
		}

		render.Status(r, statusCode)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("user info edited successfully", slog.String("user_id", userID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "user info updated successfully"})
}

func (h *userHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.GetUserInfo"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		log.Error("user_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "user_id is required"})
		return
	}

	log.Info("fetching user info", slog.String("user_id", userID))

	userInfo, err := h.userService.GetUserInfo(r.Context(), userID)
	if err != nil {
		var statusCode int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
			log.Warn("user not found", slog.String("user_id", userID))
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("failed to fetch user info", slog.String("error", err.Error()))
		}
		render.Status(r, statusCode)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("user info fetched successfully", slog.String("user_id", userID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, userInfo)
}

func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.GetUsers"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	usersInfo, err := h.userService.GetUsersInfo(r.Context())
	if err != nil {
		log.Error("failed to fetching users info", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("users info fetched successfully")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, usersInfo)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.DeleteUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID := chi.URLParam(r, "user_id")

	if userID == "" {
		log.Error("user_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "user_id is required"})
		return
	}

	log.Info("attempting to delete user", slog.String("user_id", userID))

	err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
			log.Warn("user not found", slog.String("user_id", userID))
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("failed to delete user", slog.String("user_id", userID), slog.String("error", err.Error()))
		}

		render.Status(r, statusCode)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("user deleted successfully", slog.String("user_id", userID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "user deleted successfully"})
}

func (h *userHandler) GetUserInfoByToken(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.user.GetUserInfoByToken"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		log.Error("failed to get userID from context")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	log.Info("fetching user info by token", slog.String("user_id", userID))

	userInfo, err := h.userService.GetUserInfo(r.Context(), userID)
	if err != nil {
		var statusCode int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
			log.Warn("user not found", slog.String("user_id", userID))
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("failed to fetch user info", slog.String("error", err.Error()))
		}

		render.Status(r, statusCode)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("user info fetched successfully", slog.String("user_id", userID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, userInfo)
}
