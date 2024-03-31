package repositories

import (
	groupsModels "sad/internal/models/groups"
	subjectsModels "sad/internal/models/subjects"
	usersModels "sad/internal/models/users"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	GetById(c *fiber.Ctx, userId string) (*usersModels.UserCredentials, error)
	GetByLogin(c *fiber.Ctx, login string) (*usersModels.UserRepoModel, error)
	Create(c *fiber.Ctx, user usersModels.User) error
	ChangeUserRole(c *fiber.Ctx, userId string, newRole usersModels.UserRole) error
	CheckUserExists(c *fiber.Ctx, userId string) (bool, error)
	GetUsersInfoByIds(c *fiber.Ctx, usersId []string) ([]usersModels.UserInfoRepoModel, error)
}

type GroupsRepository interface {
	Create(c *fiber.Ctx, group groupsModels.Group) error
	GetAll(c *fiber.Ctx) ([]groupsModels.GroupRepoModel, error)
	GetById(c *fiber.Ctx, groupId string) (*groupsModels.GroupRepoModel, error)
	AddUserToGroup(c *fiber.Ctx, groupId string, userId string) error
	DeleteUserFromGroup(c *fiber.Ctx, groupId string, userId string) error
	IsUserInGroup(c *fiber.Ctx, groupId, userId string) (bool, error)
	GetByIdWithUsers(c *fiber.Ctx, groupId string) (*groupsModels.GroupWithUsersRepo, error)
	DeleteGroup(c *fiber.Ctx, groupId string) error
	UpdateGroup(c *fiber.Ctx, group groupsModels.Group) error
	CheckGroupExists(c *fiber.Ctx, groupId string) (bool, error)
}

type SubjectsRepository interface {
	Create(c *fiber.Ctx, subject subjectsModels.Subject) error
	GetAll(c *fiber.Ctx) ([]subjectsModels.Subject, error)
	GetById(c *fiber.Ctx, groupId string) (*subjectsModels.Subject, error)
	DeleteSubject(c *fiber.Ctx, subjectId string) error
	AddSubjectToGroup(c *fiber.Ctx, subjectId string, groupId string) error
	DeleteSubjectFromGroup(c *fiber.Ctx, subjectId string, groupId string) error
	UpdateSubject(c *fiber.Ctx, subject subjectsModels.Subject) error
	IsSubjectInGroup(c *fiber.Ctx, subjectId, groupId string) (bool, error)
}
