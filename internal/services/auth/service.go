package auth

import (
	"errors"
	"log"
	"sad/internal/config"
	errorsModels "sad/internal/models/errors"
	userModels "sad/internal/models/user"
	"time"

	authModels "sad/internal/models/auth"
	"sad/internal/repositories"
	"sad/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/gofiber/fiber/v2"
)

type service struct {
	userRepository repositories.UserRepository
}

func NewService(userRepository repositories.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) Login(c *fiber.Ctx, user userModels.UserCredentials) (string, error) {
	log.Printf("Attempting login for user: %s", user.Login)
	existedUser, err := s.userRepository.GetByLogin(c, user.Login)

	if err != nil {
		log.Printf("Repo get user by login error: %s", err.Error())
		return "", err
	}

	if existedUser == nil {
		log.Printf("User not found for login: %s", user.Login)
		return "", errorsModels.ErrUserNotFound
	}

	if !utils.CompareHashPassword(user.Password, existedUser.Password) {
		log.Printf("Invalid credentials for user: %s", user.Login)
		return "", errorsModels.ErrInvalidCredentials
	}

	token, err := s.GetJWTToken(existedUser.UUID)

	if err != nil {
		log.Printf("Error generating JWT token for user: %s, error: %s", user.Login, err.Error())
		return "", errorsModels.ErrServer
	}

	log.Printf("User logged in successfully: %s", user.Login)
	return token, nil
}

func (s *service) Register(c *fiber.Ctx, user authModels.UserRegistrationRequest) (string, error) {
	log.Printf("Attempting to register new user: %s", user.Login)

	if user.Login == "" {
		log.Printf("Validation error: login is required")
		return "", errors.New("login is required")
	}
	if user.Name == "" {
		log.Printf("Validation error: name is required")
		return "", errors.New("name is required")
	}
	if user.Password == "" {
		log.Printf("Validation error: password is required")
		return "", errors.New("password is required")
	}
	if len(user.Password) < 8 {
		log.Printf("Validation error: password length must be more or equal 8 symbols")
		return "", errors.New("password length must be more or equal 8 symbols")
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password for user '%s': %s", user.Login, err.Error())
		return "", errorsModels.ErrServer
	}

	newUser := userModels.User{
		UUID:     uuid.New().String(),
		Name:     user.Name,
		Login:    user.Login,
		Password: hashedPassword,
		Role:     userModels.Student,
	}

	if err := s.userRepository.Create(c, newUser); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("User with login '%s' already exists", user.Login)
				return "", errorsModels.ErrUserExists
			default:
				log.Printf("Error creating user '%s' in the repository: %s", newUser.Login, err.Error())
				return "", errorsModels.ErrServer
			}
		} else {
			log.Printf("Error creating user '%s' in the repository: %s", newUser.Login, err.Error())
			return "", errorsModels.ErrServer
		}
	}

	log.Printf("User '%s' registered successfully with UUID: %s", newUser.Login, newUser.UUID)
	return newUser.UUID, nil
}

func (s *service) GetJWTToken(uuid string) (string, error) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("Error loading config: %v", err)
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
		log.Printf("Error signing the token: %v", err)
		return "", err
	}

	log.Printf("JWT token generated for UUID: %s", uuid)
	return tokenString, nil
}
