package mocks

//
//import (
//	"database/sql"
//	"github.com/gofiber/fiber/v2"
//	"github.com/stretchr/testify/assert"
//	usersModels "sad/internal/models/users"
//	"sad/internal/services/auth"
//	mocks "sad/test/repositories/users"
//	"testing"
//)
//
//func TestLogin(t *testing.T) {
//	mockUserRepo := new(mocks.UsersRepository)
//	service := auth.NewService(mockUserRepo)
//
//	ctx := &fiber.Ctx{}
//	user := usersModels.UserCredentials{
//		Login:    "testuser",
//		Password: "testpassword",
//	}
//
//	existedUser := &usersModels.UserRepoModel{
//		Id:       sql.NullString{String: "user-id", Valid: true},
//		Login:    sql.NullString{String: "testuser", Valid: true},
//		Password: sql.NullString{String: "$2a$14$Pr1H5rLHX00wBmst678HxuxbgSaTyQbMzTXwWt2dIB0PrIuSg6BYe", Valid: true}, // Хеш пароля "testpassword"
//	}
//
//	mockUserRepo.On("GetByLogin", ctx, user.Login).Return(existedUser, nil)
//
//	token, err := service.Login(ctx, user)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, token)
//}
