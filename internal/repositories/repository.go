package repositories

import (
	"context"
	gradesModels "sad/internal/models/grades"
	groupsModels "sad/internal/models/groups"
	subjectsModels "sad/internal/models/subjects"
	usersModels "sad/internal/models/users"
)

type UserRepository interface {
	GetById(ctx context.Context, userId string) (*usersModels.UserCredentials, error)
	GetByLogin(ctx context.Context, login string) (*usersModels.UserRepoModel, error)
	Create(ctx context.Context, user usersModels.User) error
	ChangeUserInfo(ctx context.Context, userId string, newRole usersModels.UserRole, newName string) error
	CheckUserExists(ctx context.Context, userId string) (bool, error)
	GetUserInfo(ctx context.Context, userId string) (*usersModels.UserInfoRepoModel, error)
	GetUsersInfo(ctx context.Context) ([]usersModels.UserInfoRepoModel, error)
	DeleteUser(ctx context.Context, userId string) error
	GetAvailableTeachers(ctx context.Context, teacherName string) ([]usersModels.UserInfoRepoModel, error)
}

type GroupsRepository interface {
	Create(ctx context.Context, group groupsModels.Group) error
	GetAll(ctx context.Context) ([]groupsModels.GroupRepoModel, error)
	GetById(ctx context.Context, groupId string) (*groupsModels.GroupRepoModel, error)
	AddUserToGroup(ctx context.Context, groupId string, userId string) error
	DeleteUserFromGroup(ctx context.Context, groupId string, userId string) error
	IsUserInGroup(ctx context.Context, groupId, userId string) (bool, error)
	GetWithDetailsById(ctx context.Context, groupId string) (*groupsModels.GroupDetailsRepo, error)
	DeleteGroup(ctx context.Context, groupId string) error
	UpdateGroup(ctx context.Context, group groupsModels.Group) error
	CheckGroupExists(ctx context.Context, groupId string) (bool, error)
	GetAvailableNewUsers(ctx context.Context, groupId, login string) ([]usersModels.UserInfoRepoModel, error)
	GetGroupsWithSubjectsByTeacher(ctx context.Context, teacherId string) ([]subjectsModels.GroupsWithSubjectsRepoModel, error)
	GetGroupsBySubjectAndTeacher(ctx context.Context, teacherId, subjectId string) ([]groupsModels.GroupRepoModel, error)
}

type SubjectsRepository interface {
	Create(ctx context.Context, subject subjectsModels.Subject) (*subjectsModels.Subject, error)
	GetAll(ctx context.Context) ([]subjectsModels.SubjectRepoModel, error)
	GetById(ctx context.Context, groupId string) (*subjectsModels.SubjectRepoModel, error)
	DeleteSubject(ctx context.Context, subjectId string) error
	AddSubjectToGroup(ctx context.Context, subjectGroup subjectsModels.SubjectGroup) error
	DeleteSubjectFromGroup(ctx context.Context, subjectId string, groupId string) error
	UpdateSubject(ctx context.Context, subject subjectsModels.Subject) error
	IsSubjectInGroup(ctx context.Context, subjectId, groupId string) (bool, error)
	GetSubjectTeacherId(ctx context.Context, subjectId, teacherId string) (string, error)
	AddTeacherToSubject(ctx context.Context, subjectId, teacherId string) error
	GetByIdWithDetails(ctx context.Context, subjectId string) (*subjectsModels.SubjectRepoModel, error)
	GetSubjectsByTeacherId(ctx context.Context, teacherId string) ([]subjectsModels.SubjectRepoModel, error)
	GetNewSubjectsForTeacher(ctx context.Context, teacherId string) ([]subjectsModels.SubjectRepoModel, error)
	UpdateTeacherSubjects(ctx context.Context, teacherId string, subjects []subjectsModels.Subject) error
}

type GradesRepository interface {
	Create(ctx context.Context, grade gradesModels.Grade) error
	Delete(ctx context.Context, gradeId string) error
	Update(ctx context.Context, grade gradesModels.Grade) error
	GetAllStudentGrades(ctx context.Context, userId string, isFinal bool, subjectId *string) ([]gradesModels.GradeInfoRepoModel, error)
	GetById(ctx context.Context, gradeId string) (*gradesModels.GradeRepoModel, error)
	GetStudentsGradesBySubjectAndGroup(ctx context.Context, subjectId, studentId string, isFinal *bool) ([]gradesModels.UserSubjectGradesRepoModel, error)
	GetAllGradesInfo(ctx context.Context) ([]gradesModels.GradesReportRecordRepoModel, error)
	//GetGroupGradesBySubjectId(ctx context.Context, subjectId string, groupId string) ([]gradesModels.GradeRepoModel, error)
}
