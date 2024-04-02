package services

import (
	"sad/internal/models/auth"

	"sad/internal/models/users"

	"sad/internal/models/groups"

	"sad/internal/models/subjects"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Login(c *fiber.Ctx, user usersModels.UserCredentials) (string, error)
	Register(c *fiber.Ctx, user authModels.UserRegistrationRequest) (string, error)
}

type UserService interface {
	EditRole(c *fiber.Ctx, userId string, newRole usersModels.UserRole) error
	CheckUserIsAdmin(c *fiber.Ctx) (bool, error)
}

type GroupsService interface {
	Create(c *fiber.Ctx, number string) error
	GetAll(c *fiber.Ctx) ([]groupsModels.Group, error)
	GetById(c *fiber.Ctx, groupId string) (*groupsModels.Group, error)
	GetByIdWithUsers(c *fiber.Ctx, groupId string) (*groupsModels.GroupWithUsers, error)
	DeleteGroup(c *fiber.Ctx, groupId string) error
	AddUserToGroup(c *fiber.Ctx, groupId string, userId string) error
	DeleteUserFromGroup(c *fiber.Ctx, groupId string, userId string) error
	UpdateGroup(c *fiber.Ctx, groupId string, group groupsModels.Group) error
}

type SubjectsService interface {
	Create(c *fiber.Ctx, name string) error
	GetAll(c *fiber.Ctx) ([]subjectsModels.Subject, error)
	DeleteSubject(c *fiber.Ctx, subjectId string) error
	AddSubjectToGroup(c *fiber.Ctx, subjectId string, groupId string) error
	DeleteSubjectFromGroup(c *fiber.Ctx, subjectId string, groupId string) error
	UpdateSubject(c *fiber.Ctx, subjectId string, subject subjectsModels.Subject) error
}
