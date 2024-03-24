package app

import (
	"sad/internal/config"
	database "sad/internal/db"
	"sad/internal/handlers/auth"
	"sad/internal/handlers/user"
	"sad/internal/repositories"
	userRepository "sad/internal/repositories/user"
	"sad/internal/services"
	authService "sad/internal/services/auth"
	userService "sad/internal/services/user"

	"github.com/jackc/pgx/v5"
)

type serviceProvider struct {
	authService services.AuthService

	userService services.UserService

	userRepository repositories.UserRepository

	authHandler auth.AuthHandler

	userHandler user.UserHandler

	db *pgx.Conn
}

func newServiceProvider(config config.Config) (*serviceProvider, error) {
	db, err := database.NewDBConnection(config)

	if err != nil {
		return nil, err
	}

	return &serviceProvider{
		db: db,
	}, nil
}

func (s *serviceProvider) UserRepository() repositories.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.db)
	}

	return s.userRepository
}

func (s *serviceProvider) AuthService() services.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.UserRepository())
	}

	return s.authService
}

func (s *serviceProvider) UserService() services.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository())
	}

	return s.userService
}

func (s *serviceProvider) AuthHandler() auth.AuthHandler {
	if s.authHandler == nil {
		s.authHandler = auth.NewAuthHandler(s.AuthService())
	}

	return s.authHandler
}

func (s *serviceProvider) UserHandler() user.UserHandler {
	if s.userHandler == nil {
		s.userHandler = user.NewUserHandler(s.UserService())
	}

	return s.userHandler
}
