package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sad/internal/config"
	authModels "sad/internal/models/auth"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(config config.Config, logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.auth.middleware"

			log := logger.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			var tokenString string
			authorization := r.Header.Get("Authorization")

			if strings.HasPrefix(authorization, "Bearer ") {
				tokenString = strings.TrimPrefix(authorization, "Bearer ")
			}

			if tokenString == "" {
				log.Warn("Authorization header is missing or does not contain Bearer token")

				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"message": "You are not logged in"}`))
				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &authModels.Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					errMsg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])

					log.Warn(errMsg)
					return nil, errors.New(errMsg)
				}

				// Возвращаем секретный ключ для проверки подписи
				return []byte(config.JwtSecret), nil
			})

			if err != nil {
				log.Warn("Error parsing token: %v", slog.String("error", err.Error()))

				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"message": "unauthorized"}`))
				return
			}

			if claims, ok := token.Claims.(*authModels.Claims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "user_id", claims.Subject)
				log.Info("Token is valid", slog.String("userID", claims.Subject))
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				log.Error("Token is not valid or claims are not of expected type")

				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"error": "unauthorized"}`))
			}
		})

	}
}
