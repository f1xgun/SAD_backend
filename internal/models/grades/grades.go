package gradesModels

import (
	"database/sql"
	usersModels "sad/internal/models/users"
	"strconv"
	"time"
)

type Grade struct {
	Id         string    `json:"id"`
	SubjectId  string    `json:"subject_id"`
	StudentId  string    `json:"student_id"`
	TeacherId  string    `json:"teacher_id"`
	Evaluation int       `json:"evaluation"`
	CreatedAt  time.Time `json:"created_at"`
	IsFinal    *bool     `json:"is_final"`
	Comment    *string   `json:"comment"`
}

type GradeRepoModel struct {
	Id         sql.NullString
	SubjectId  sql.NullString
	StudentId  sql.NullString
	TeacherId  sql.NullString
	Evaluation sql.NullInt16
	CreatedAt  sql.NullTime
	IsFinal    sql.NullBool
	Comment    sql.NullString
}

type GradeInfo struct {
	Id          string    `json:"id"`
	SubjectName string    `json:"subject_name,omitempty"`
	TeacherName string    `json:"teacher_name,omitempty"`
	Evaluation  int       `json:"evaluation"`
	CreatedAt   time.Time `json:"created_at"`
	IsFinal     *bool     `json:"is_final,omitempty"`
	Comment     *string   `json:"comment,omitempty"`
}

type GradeInfoRepoModel struct {
	Id          sql.NullString
	SubjectName sql.NullString
	TeacherName sql.NullString
	Evaluation  sql.NullInt16
	CreatedAt   sql.NullTime
	IsFinal     sql.NullBool
	Comment     sql.NullString
}

type UserSubjectGrades struct {
	Student usersModels.UserInfo `json:"student"`
	Grades  []GradeInfo          `json:"grades"`
}

type UserSubjectGradesRepoModel struct {
	Student usersModels.UserInfoRepoModel
	Grades  []GradeInfoRepoModel
}

type GradesReportRecordRepoModel struct {
	Student     usersModels.UserInfoRepoModel
	Teacher     usersModels.UserInfoRepoModel
	GroupNumber sql.NullString
	GradeInfo   GradeInfoRepoModel
}

func (m *GradesReportRecordRepoModel) ToCsvString() []string {
	var studentMiddleName string
	var teacherMiddleName string
	var gradeComment string
	gradeIsFinal := "Нет"

	if m.Student.MiddleName.Valid {
		studentMiddleName = m.Student.MiddleName.String
	}

	if m.Teacher.MiddleName.Valid {
		teacherMiddleName = m.Teacher.MiddleName.String
	}

	if m.GradeInfo.Comment.Valid {
		gradeComment = m.GradeInfo.Comment.String
	}

	if m.GradeInfo.IsFinal.Valid && m.GradeInfo.IsFinal.Bool {
		gradeIsFinal = "Да"
	}

	return []string{
		m.Student.LastName.String,
		m.Student.Name.String,
		studentMiddleName,
		m.GroupNumber.String,
		m.GradeInfo.SubjectName.String,
		strconv.Itoa(int(m.GradeInfo.Evaluation.Int16)),
		m.GradeInfo.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		gradeComment,
		gradeIsFinal,
		m.Teacher.LastName.String,
		m.Teacher.Name.String,
		teacherMiddleName,
	}
}
