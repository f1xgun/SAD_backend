package userModels

type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	UUID     string
	Name     string
	Login    string
	Password string
	Role     string
}
