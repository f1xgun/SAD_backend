package services

import (
	"context"
	authModels "sad/internal/models/auth"
	gradesModels "sad/internal/models/grades"

	usersModels "sad/internal/models/users"

	groupsModels "sad/internal/models/groups"

	subjectsModels "sad/internal/models/subjects"
)

type AuthService interface {
	Login(ctx context.Context, user usersModels.UserCredentials) (string, error)
	Register(ctx context.Context, user authModels.UserRegistrationRequest) error
}

type UserService interface {
	EditUser(ctx context.Context, userId string, newRole usersModels.UserRole, newName string) error
	CheckIsUserRoleAllowed(ctx context.Context, allowedRoles []usersModels.UserRole, userId string) (bool, error)
	GetUserInfo(ctx context.Context, userId string) (*usersModels.UserInfo, error)
	GetUsersInfo(ctx context.Context) ([]usersModels.UserInfo, error)
	DeleteUser(ctx context.Context, userId string) error
}

type GroupsService interface {
	Create(ctx context.Context, number string) error
	GetAll(ctx context.Context) ([]groupsModels.Group, error)
	GetById(ctx context.Context, groupId string) (*groupsModels.Group, error)
	GetWithDetailsById(ctx context.Context, groupId string) (*groupsModels.GroupDetails, error)
	DeleteGroup(ctx context.Context, groupId string) error
	AddUserToGroup(ctx context.Context, groupId string, userId string) error
	DeleteUserFromGroup(ctx context.Context, groupId string, userId string) error
	UpdateGroup(ctx context.Context, groupId string, group groupsModels.Group) error
	GetAvailableNewUsers(ctx context.Context, groupId, login string) ([]usersModels.UserInfo, error)
	GetGroupsWithSubjectsByTeacher(ctx context.Context, teacherId string) ([]subjectsModels.GroupsWithSubjects, error)
	GetGroupsBySubjectAndTeacher(ctx context.Context, teacherId, subjectId string) ([]groupsModels.Group, error)
}

type SubjectsService interface {
	Create(ctx context.Context, name string) error
	GetAll(ctx context.Context) ([]subjectsModels.Subject, error)
	DeleteSubject(ctx context.Context, subjectId string) error
	AddSubjectToGroup(ctx context.Context, subjectGroup subjectsModels.SubjectGroup) error
	DeleteSubjectFromGroup(ctx context.Context, subjectId string, groupId string) error
	UpdateSubject(ctx context.Context, subjectId string, subject subjectsModels.Subject) error
	GetAvailableTeachers(ctx context.Context, teacherName string) ([]usersModels.UserInfo, error)
	GetByIdWithDetails(ctx context.Context, subjectId string) (*subjectsModels.Subject, error)
	GetSubjectsByTeacherId(ctx context.Context, teacherId string) ([]subjectsModels.Subject, error)
	GetNewAvailableSubjectsForTeacher(ctx context.Context, teacherId string) ([]subjectsModels.Subject, error)
	EditTeacherSubjects(ctx context.Context, teacherId string, subjects []subjectsModels.Subject) error
}

type GradesService interface {
	Create(ctx context.Context, grade gradesModels.Grade) error
	Delete(ctx context.Context, gradeId string) error
	Update(ctx context.Context, gradeId string, evaluation *int, comment *string) error
	GetAllStudentGrades(ctx context.Context, userId string, isFinal bool, subjectId *string) ([]gradesModels.GradeInfo, error)
	GetStudentsGradesBySubjectAndGroup(ctx context.Context, subjectId, groupId string, isFinal *bool) ([]gradesModels.UserSubjectGrades, error)
	GetGradesInCsv(ctx context.Context) (string, error)
}
