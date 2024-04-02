package app

import (
	"sad/internal/config"
	"sad/internal/db"
	"sad/internal/handlers/auth"
	"sad/internal/handlers/groups"
	"sad/internal/handlers/subjects"
	"sad/internal/handlers/user"
	"sad/internal/repositories"
	groupsRepository "sad/internal/repositories/groups"
	subjectsRepository "sad/internal/repositories/subjects"
	userRepository "sad/internal/repositories/user"
	"sad/internal/services"
	authService "sad/internal/services/auth"
	groupsService "sad/internal/services/groups"
	subjectsService "sad/internal/services/subjects"
	userService "sad/internal/services/user"

	"github.com/jackc/pgx/v5"
)

type serviceProvider struct {
	authService services.AuthService

	userService services.UserService

	groupsService services.GroupsService

	subjectsService services.SubjectsService

	userRepository repositories.UserRepository

	groupsRepository repositories.GroupsRepository

	subjectsRepository repositories.SubjectsRepository

	authHandler auth.AuthHandler

	userHandler users.UserHandler

	groupsHandler groups.GroupsHandler

	subjectsHandler subjects.SubjectsHandler

	db *pgx.Conn
}

func newServiceProvider(config config.Config) (*serviceProvider, error) {
	database, err := db.NewDBConnection(config)

	if err != nil {
		return nil, err
	}

	return &serviceProvider{
		db: database,
	}, nil
}

func (s *serviceProvider) UserRepository() repositories.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.db)
	}

	return s.userRepository
}

func (s *serviceProvider) GroupsRepository() repositories.GroupsRepository {
	if s.groupsRepository == nil {
		s.groupsRepository = groupsRepository.NewRepository(s.db)
	}

	return s.groupsRepository
}

func (s *serviceProvider) SubjectsRepository() repositories.SubjectsRepository {
	if s.subjectsRepository == nil {
		s.subjectsRepository = subjectsRepository.NewRepository(s.db)
	}

	return s.subjectsRepository
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

func (s *serviceProvider) GroupsService() services.GroupsService {
	if s.groupsService == nil {
		s.groupsService = groupsService.NewService(s.GroupsRepository(), s.UserRepository())
	}

	return s.groupsService
}

func (s *serviceProvider) SubjectsService() services.SubjectsService {
	if s.subjectsService == nil {
		s.subjectsService = subjectsService.NewService(s.GroupsRepository(), s.SubjectsRepository())
	}

	return s.subjectsService
}

func (s *serviceProvider) NewAuthHandler() auth.AuthHandler {
	if s.authHandler == nil {
		s.authHandler = auth.NewAuthHandler(s.AuthService())
	}

	return s.authHandler
}

func (s *serviceProvider) NewUserHandler() users.UserHandler {
	if s.userHandler == nil {
		s.userHandler = users.NewUserHandler(s.UserService())
	}

	return s.userHandler
}

func (s *serviceProvider) NewGroupsHandler() groups.GroupsHandler {
	if s.groupsHandler == nil {
		s.groupsHandler = groups.NewGroupsHandler(s.GroupsService())
	}

	return s.groupsHandler
}

func (s *serviceProvider) NewSubjectsHandler() subjects.SubjectsHandler {
	if s.subjectsHandler == nil {
		s.subjectsHandler = subjects.NewSubjectsHandler(s.SubjectsService())
	}

	return s.subjectsHandler
}
