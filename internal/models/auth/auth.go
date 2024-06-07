package authModels

import "github.com/golang-jwt/jwt/v5"

type UserRegistrationRequest struct {
	LastName   string  `json:"last_name"`
	Name       string  `json:"name"`
	MiddleName *string `json:"middle_name"`
	Login      string  `json:"login"`
	Password   string  `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
}
