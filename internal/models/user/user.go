package userModels

type UserCredentials struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

type User struct {
	UUID     string
	Name     string
	Login    string
	Password string
	Role     UserRole
}

type UserRepoModel struct {
	UUID     string
	Login    string
	Password string
}

type UserRole string

const (
	Student UserRole = "student"
	Teacher UserRole = "teacher"
	Admin   UserRole = "admin"
)
