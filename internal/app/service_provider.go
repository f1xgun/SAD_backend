package app

import (
	"sad/internal/config"
	database "sad/internal/db"
	"sad/internal/handlers/auth"
	"sad/internal/repositories"
	authRepository "sad/internal/repositories/auth"
	"sad/internal/services"
	authService "sad/internal/services/auth"

	"github.com/jackc/pgx/v5"
)

type serviceProvider struct {
	authService services.AuthService

	authRepository repositories.AuthRepository

	authHandler auth.AuthHandler

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

func (s *serviceProvider) AuthRepository() repositories.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.db)
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService() services.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository())
	}

	return s.authService
}

func (s *serviceProvider) AuthHandler() auth.AuthHandler {
	if s.authHandler == nil {
		s.authHandler = auth.NewAuthHandler(s.AuthService())
	}

	return s.authHandler
}
