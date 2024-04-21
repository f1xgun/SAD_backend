package gradesModels

import (
	"database/sql"
	"time"
)

type Grade struct {
	Id         string    `json:"id"`
	SubjectId  string    `json:"subject_id"`
	StudentId  string    `json:"student_id"`
	TeacherId  string    `json:"teacher_id"`
	Evaluation int       `json:"evaluation"`
	CreatedAt  time.Time `json:"created_at"`
}

type GradeRepoModel struct {
	Id         sql.NullString
	SubjectId  sql.NullString
	StudentId  sql.NullString
	TeacherId  sql.NullString
	Evaluation sql.NullInt16
	CreatedAt  sql.NullTime
}

type GradeInfo struct {
	Id          string    `json:"id"`
	SubjectName string    `json:"subject_name"`
	TeacherName string    `json:"teacher_name"`
	Evaluation  int       `json:"evaluation"`
	CreatedAt   time.Time `json:"created_at"`
}

type GradeInfoRepoModel struct {
	Id          sql.NullString
	SubjectName sql.NullString
	TeacherName sql.NullString
	Evaluation  sql.NullInt16
	CreatedAt   sql.NullTime
}
