package service

import (
	"sad/internal/models/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Login(c *fiber.Ctx, user *auth.UserLoginRequest) (string, error)
	Register(c *fiber.Ctx, user *auth.UserRegistrationRequest) (string, error)
}
