package mocks

import (
	authModels "sad/internal/models/auth"
	"sad/internal/models/users"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/mock"
)

type AuthService struct {
	mock.Mock
}

func (m *AuthService) Login(ctx *fiber.Ctx, user users.UserCredentials) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *AuthService) Register(ctx *fiber.Ctx, user authModels.UserRegistrationRequest) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}
