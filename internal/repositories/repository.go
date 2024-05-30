package repositories

import (
	gradesModels "sad/internal/models/grades"
	groupsModels "sad/internal/models/groups"
	subjectsModels "sad/internal/models/subjects"
	usersModels "sad/internal/models/users"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	GetById(c *fiber.Ctx, userId string) (*usersModels.UserCredentials, error)
	GetByLogin(c *fiber.Ctx, login string) (*usersModels.UserRepoModel, error)
	Create(c *fiber.Ctx, user usersModels.User) error
	ChangeUserInfo(c *fiber.Ctx, userId string, newRole usersModels.UserRole, newName string) error
	CheckUserExists(c *fiber.Ctx, userId string) (bool, error)
	GetUserInfo(c *fiber.Ctx, userId string) (*usersModels.UserInfoRepoModel, error)
	GetUsersInfo(c *fiber.Ctx) ([]usersModels.UserInfoRepoModel, error)
	DeleteUser(c *fiber.Ctx, userId string) error
	GetAvailableTeachers(c *fiber.Ctx, teacherName string) ([]usersModels.UserInfoRepoModel, error)
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
	GetAvailableNewUsers(c *fiber.Ctx, groupId, login string) ([]usersModels.UserInfoRepoModel, error)
	GetGroupsWithSubjectsByTeacher(c *fiber.Ctx, teacherId string) ([]subjectsModels.GroupsWithSubjectsRepoModel, error)
	GetGroupsBySubjectAndTeacher(c *fiber.Ctx, teacherId, subjectId string) ([]groupsModels.GroupRepoModel, error)
}

type SubjectsRepository interface {
	Create(c *fiber.Ctx, subject subjectsModels.Subject) error
	GetAll(c *fiber.Ctx) ([]subjectsModels.SubjectRepoModel, error)
	GetById(c *fiber.Ctx, groupId string) (*subjectsModels.SubjectRepoModel, error)
	DeleteSubject(c *fiber.Ctx, subjectId string) error
	AddSubjectToGroup(c *fiber.Ctx, subjectGroup subjectsModels.SubjectGroup) error
	DeleteSubjectFromGroup(c *fiber.Ctx, subjectId string, groupId string) error
	UpdateSubject(c *fiber.Ctx, subject subjectsModels.Subject) error
	IsSubjectInGroup(c *fiber.Ctx, subjectId, groupId string) (bool, error)
	GetSubjectTeacherId(c *fiber.Ctx, subjectId, teacherId string) (string, error)
	AddTeacherToSubject(c *fiber.Ctx, subjectTeacherId, subjectId, teacherId string) error
	GetByIdWithDetails(c *fiber.Ctx, subjectId string) (*subjectsModels.SubjectInfoRepoModel, error)
	GetSubjectsByTeacherId(c *fiber.Ctx, teacherId string) ([]subjectsModels.SubjectRepoModel, error)
}

type GradesRepository interface {
	Create(c *fiber.Ctx, grade gradesModels.Grade) error
	Delete(c *fiber.Ctx, gradeId string) error
	Update(c *fiber.Ctx, grade gradesModels.Grade) error
	GetAllStudentGrades(c *fiber.Ctx, userId string, isFinal bool, subjectId *string) ([]gradesModels.GradeInfoRepoModel, error)
	GetById(c *fiber.Ctx, gradeId string) (*gradesModels.GradeRepoModel, error)
	GetStudentsGradesBySubjectAndGroup(c *fiber.Ctx, subjectId, studentId string, isFinal *bool) ([]gradesModels.UserSubjectGradesRepoModel, error)
	//GetGroupGradesBySubjectId(c *fiber.Ctx, subjectId string, groupId string) ([]gradesModels.GradeRepoModel, error)
}
