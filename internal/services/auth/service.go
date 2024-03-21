package auth

import (
	"sad/internal/config"
	"sad/internal/models/errors"
	userModels "sad/internal/models/user"
	"time"

	authModels "sad/internal/models/auth"
	"sad/internal/repositories"
	"sad/internal/utils"

	"github.com/golang-jwt/jwt/v5"
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
		return "", errors.ErrUserNotFound
	}

	if !utils.CompareHashPassword(user.Password, existedUser.Password) {
		return "", errors.ErrInvalidCredentials
	}

	token, err := s.GetJWTToken(existedUser.UUID)

	if err != nil {
		return "", errors.ErrServer
	}

	return token, nil
}

func (s *service) Register(c *fiber.Ctx, user authModels.UserRegistrationRequest) (string, error) {
	existedUser, err := s.authRepository.GetByLogin(c, user.Login)

	if err != nil {
		return "", err
	}

	if existedUser != nil {
		return "", errors.ErrUserExists
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)

	if err != nil {
		return "", errors.ErrServer
	}

	newUser := userModels.User{
		UUID:     uuid.New().String(),
		Name:     user.Name,
		Login:    user.Login,
		Password: hashedPassword,
		Role:     "default",
	}

	if err := s.authRepository.Create(c, newUser); err != nil {
		return "", errors.ErrServer
	}

	return newUser.UUID, nil
}

func (s *service) GetJWTToken(uuid string) (string, error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	claims := authModels.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid,
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(config.JwtExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	tokenByte := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
