package services

import (
	"sad/internal/models/auth"

	userModels "sad/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Login(c *fiber.Ctx, user userModels.UserCredentials) (string, error)
	Register(c *fiber.Ctx, user auth.UserRegistrationRequest) (string, error)
}
