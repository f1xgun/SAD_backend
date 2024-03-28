package repositories

import (
	groupsModels "sad/internal/models/groups"
	usersModels "sad/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	GetById(c *fiber.Ctx, userId string) (*usersModels.UserCredentials, error)
	GetByLogin(c *fiber.Ctx, login string) (*usersModels.UserRepoModel, error)
	Create(c *fiber.Ctx, user usersModels.User) error
	ChangeUserRole(c *fiber.Ctx, userId string, newRole usersModels.UserRole) error
	CheckUserExists(c *fiber.Ctx, userId string) (bool, error)
	GetUsersInfoByIds(c *fiber.Ctx, usersId []string) ([]usersModels.UserInfo, error)
}

type GroupsRepository interface {
	Create(c *fiber.Ctx, group groupsModels.Group) error
	GetAll(c *fiber.Ctx) ([]groupsModels.Group, error)
	GetById(c *fiber.Ctx, groupId string) (*groupsModels.Group, error)
	AddUserToGroup(c *fiber.Ctx, groupId string, userId string) error
	DeleteUserFromGroup(c *fiber.Ctx, groupId string, userId string) error
	IsUserInGroup(c *fiber.Ctx, groupId, userId string) (bool, error)
	GetGroupUsers(c *fiber.Ctx, groupId string) ([]string, error)
	DeleteGroup(c *fiber.Ctx, groupId string) error
	UpdateGroup(c *fiber.Ctx, group groupsModels.Group) error
}
