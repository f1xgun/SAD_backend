package users

import (
	"encoding/json"
	"log/slog"
	"net/http"
	usersModels "sad/internal/models/users"
	"sad/internal/services"

	"github.com/go-chi/chi/v5/middleware"
)

func AllowedRoleMiddleware(usersService services.UserService, allowedRoles []usersModels.UserRole, logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.users.AllowedRoleMiddleware"

			log := logger.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			log.Info("Checking if user has allowed role")

			userID, ok := r.Context().Value("user_id").(string)
			if !ok {
				log.Error("Failed to assert type for userID from context")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "failed to get user id"})
				return
			}

			userHasAllowedRole, err := usersService.CheckIsUserRoleAllowed(r.Context(), allowedRoles, userID)
			if err != nil {
				log.Error("Error checking if user has allowed role", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			if !userHasAllowedRole {
				log.Warn("User does not have permission", slog.String("userID", userID))
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"message": "no permission"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
