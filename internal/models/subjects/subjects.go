package subjectsModels

import (
	"database/sql"
	groupsModels "sad/internal/models/groups"
	usersModels "sad/internal/models/users"
)

type Subject struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	TeacherId string `json:"teacher_id,omitempty"`
}

type SubjectInfo struct {
	Id      string               `json:"id"`
	Name    string               `json:"name"`
	Teacher usersModels.UserInfo `json:"teacher"`
}

type SubjectGroup struct {
	SubjectId string `json:"subject_id"`
	TeacherId string `json:"teacher_id"`
	GroupId   string `json:"group_id"`
}

type GroupsWithSubjects struct {
	Group    groupsModels.Group `json:"group"`
	Subjects []Subject          `json:"subjects"`
}

type SubjectRepoModel struct {
	Id   sql.NullString
	Name sql.NullString
}

type SubjectGroupRepoModel struct {
	SubjectId sql.NullString
	GroupId   sql.NullString
}

type SubjectInfoRepoModel struct {
	Id      sql.NullString
	Name    sql.NullString
	Teacher usersModels.UserInfoRepoModel
}

type GroupsWithSubjectsRepoModel struct {
	Group    groupsModels.GroupRepoModel
	Subjects []SubjectRepoModel
}
