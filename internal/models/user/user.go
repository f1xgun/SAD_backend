package usersModels

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
	Id       string
	Login    string
	Password string
}

type UserRole string

const (
	Student UserRole = "student"
	Teacher UserRole = "teacher"
	Admin   UserRole = "admin"
)
