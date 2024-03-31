package usersModels

import "database/sql"

type UserCredentials struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

type UserInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type User struct {
	Id       string
	Name     string
	Login    string
	Password string
	Role     UserRole
}

type UserRepoModel struct {
	Id       sql.NullString
	Login    sql.NullString
	Password sql.NullString
}

type UserInfoRepoModel struct {
	Id    sql.NullString
	Name  sql.NullString
	Login sql.NullString
}

type UserRole string

const (
	Student UserRole = "student"
	Teacher UserRole = "teacher"
	Admin   UserRole = "admin"
)
