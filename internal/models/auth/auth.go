package authModels

import "github.com/golang-jwt/jwt/v5"

type UserRegistrationRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
}
