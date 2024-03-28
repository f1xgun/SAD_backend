package services

import (
	authModels "sad/internal/models/auth"

	usersModels "sad/internal/models/user"

	groupsModels "sad/internal/models/groups"

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
	GetById(c *fiber.Ctx, groupId string) (*groupsModels.GroupWithUsers, error)
	DeleteGroup(c *fiber.Ctx, groupId string) error
	AddUserToGroup(c *fiber.Ctx, groupId string, userId string) error
	DeleteUserFromGroup(c *fiber.Ctx, groupId string, userId string) error
	UpdateGroup(c *fiber.Ctx, groupId string, group groupsModels.Group) error
}
