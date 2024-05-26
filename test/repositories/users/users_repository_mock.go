package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	usersModels "sad/internal/models/users"
)

type UsersRepository struct {
	mock.Mock
}

func (m *UsersRepository) GetById(c *fiber.Ctx, userId string) (*usersModels.UserCredentials, error) {
	args := m.Called(c, userId)
	return args.Get(0).(*usersModels.UserCredentials), args.Error(1)
}

func (m *UsersRepository) GetByLogin(c *fiber.Ctx, login string) (*usersModels.UserRepoModel, error) {
	args := m.Called(c, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModels.UserRepoModel), args.Error(1)
}

func (m *UsersRepository) Create(c *fiber.Ctx, user usersModels.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}

func (m *UsersRepository) ChangeUserInfo(c *fiber.Ctx, userId string, newRole usersModels.UserRole, newName string) error {
	args := m.Called(c, userId, newRole, newName)
	return args.Error(0)
}

func (m *UsersRepository) CheckUserExists(c *fiber.Ctx, userId string) (bool, error) {
	args := m.Called(c, userId)
	return args.Bool(0), args.Error(1)
}

func (m *UsersRepository) GetUserInfo(c *fiber.Ctx, userId string) (*usersModels.UserInfoRepoModel, error) {
	args := m.Called(c, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModels.UserInfoRepoModel), args.Error(1)
}

func (m *UsersRepository) GetUsersInfo(c *fiber.Ctx) ([]usersModels.UserInfoRepoModel, error) {
	args := m.Called(c)
	return args.Get(0).([]usersModels.UserInfoRepoModel), args.Error(1)
}

func (m *UsersRepository) DeleteUser(c *fiber.Ctx, userId string) error {
	args := m.Called(c, userId)
	return args.Error(0)
}

func (m *UsersRepository) GetAvailableTeachers(c *fiber.Ctx, teacherName string) ([]usersModels.UserInfoRepoModel, error) {
	args := m.Called(c, teacherName)
	return args.Get(0).([]usersModels.UserInfoRepoModel), args.Error(1)
}
