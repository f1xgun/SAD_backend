package services

import (
	"sad/internal/models/auth"
	gradesModels "sad/internal/models/grades"

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
	CheckIsUserRoleAllowed(c *fiber.Ctx, allowedRoles []usersModels.UserRole, userId string) (bool, error)
	GetUserInfo(c *fiber.Ctx, userId string) (*usersModels.UserInfo, error)
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
	GetAvailableNewUsers(c *fiber.Ctx, groupId, login string) ([]usersModels.UserInfo, error)
}

type SubjectsService interface {
	Create(c *fiber.Ctx, name string) error
	GetAll(c *fiber.Ctx) ([]subjectsModels.Subject, error)
	DeleteSubject(c *fiber.Ctx, subjectId string) error
	AddSubjectToGroup(c *fiber.Ctx, subjectGroup subjectsModels.SubjectGroup) error
	DeleteSubjectFromGroup(c *fiber.Ctx, subjectId string, groupId string) error
	UpdateSubject(c *fiber.Ctx, subjectId string, subject subjectsModels.Subject) error
}

type GradesService interface {
	Create(c *fiber.Ctx, grade gradesModels.Grade) error
	Delete(c *fiber.Ctx, gradeId string) error
	Update(c *fiber.Ctx, gradeId string, evaluation int) error
	GetAllStudentGrades(c *fiber.Ctx, userId string) ([]gradesModels.GradeInfo, error)
}
