package auth

import (
	userModels "sad/internal/models/user"

	"sad/internal/models/auth"
	"sad/internal/repositories"
	"sad/internal/utils"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type service struct {
	authRepository repositories.AuthRepository
}

func NewService(authRepository repositories.AuthRepository) *service {
	return &service{
		authRepository: authRepository,
	}
}

func (s *service) Login(c *fiber.Ctx, user userModels.UserCredentials) (string, error) {
	existedUser, err := s.authRepository.GetByLogin(c, user.Login)

	if err != nil {
		return "", err
	}

	if existedUser == nil {
		return "", &fiber.Error{Code: 404, Message: "User with this login does not exist"}
	}

	if !utils.CompareHashPassword(user.Password, existedUser.Password) {
		return "", &fiber.Error{Code: 401, Message: "Invalid credentials"}
	}

	return "some token", nil
}

func (s *service) Register(c *fiber.Ctx, user auth.UserRegistrationRequest) (string, error) {
	existedUser, err := s.authRepository.GetByLogin(c, user.Login)

	if err != nil {
		return "", err
	}

	if existedUser != nil {
		return "", &fiber.Error{Code: 409, Message: "User with this login already exist"}
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)

	if err != nil {
		return "", &fiber.Error{Code: 500, Message: err.Error()}
	}

	newUser := userModels.User{
		UUID:     uuid.New().String(),
		Name:     user.Name,
		Login:    user.Login,
		Password: hashedPassword,
		Role:     "default",
	}

	if err := s.authRepository.Create(c, newUser); err != nil {
		return "", &fiber.Error{Code: 500, Message: "Something went wrong"}
	}

	return newUser.UUID, nil
}
