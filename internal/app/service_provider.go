package app

import (
	"sad/internal/handlers/auth"
	"sad/internal/repositories"
	authRepository "sad/internal/repositories/auth"
	"sad/internal/services"
	authService "sad/internal/services/auth"
)

type serviceProvider struct {
	authService services.AuthService

	authRepository repositories.AuthRepository

	authHandler auth.AuthHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) AuthRepository() repositories.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository()
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
