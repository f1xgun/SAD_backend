package subjectsModels

import "database/sql"

type Subject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SubjectGroup struct {
	SubjectId string `json:"subject_id"`
	GroupId   string `json:"group_id"`
}

type SubjectRepoModel struct {
	Id   sql.NullString
	Name sql.NullString
}
