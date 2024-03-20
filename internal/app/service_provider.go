package app

import (
	"sad/internal/handlers/auth"
	"sad/internal/repository"
	authRepository "sad/internal/repository/auth"
	"sad/internal/service"
	authService "sad/internal/service/auth"
)

type serviceProvider struct {
	authService service.AuthService

	authRepository repository.AuthRepository

	authHandler auth.AuthHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) AuthRepository() repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository()
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService() service.AuthService {
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
