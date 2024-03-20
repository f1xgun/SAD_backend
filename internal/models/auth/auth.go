package auth

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRegistrationRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
