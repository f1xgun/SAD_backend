package users

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	usersMapper "sad/internal/mappers/users"
	usersModels "sad/internal/models/users"
	"testing"
)

func TestUserInfoFromRepoToService(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := usersModels.UserInfoRepoModel{
		Id:    sql.NullString{String: "user-id", Valid: true},
		Name:  sql.NullString{String: "user-name", Valid: true},
		Login: sql.NullString{String: "user-login", Valid: true},
		Role:  sql.NullString{String: "student", Valid: true},
	}

	// Вызов тестируемой функции
	entity := usersMapper.UserInfoFromRepoToService(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.Name.String, entity.Name)
	assert.Equal(t, repoModel.Login.String, entity.Login)
	assert.Equal(t, usersModels.Student, entity.Role)
}

func TestUserInfoFromRepoToService_InvalidRole(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := usersModels.UserInfoRepoModel{
		Id:    sql.NullString{String: "user-id", Valid: true},
		Name:  sql.NullString{String: "user-name", Valid: true},
		Login: sql.NullString{String: "user-login", Valid: true},
		Role:  sql.NullString{String: "invalid-role", Valid: true},
	}

	// Вызов тестируемой функции
	entity := usersMapper.UserInfoFromRepoToService(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.Name.String, entity.Name)
	assert.Equal(t, repoModel.Login.String, entity.Login)
	assert.Equal(t, usersModels.UserRole(""), entity.Role)
}

func TestUsersInfoFromRepoToService(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []usersModels.UserInfoRepoModel{
		{
			Id:    sql.NullString{String: "user-id-1", Valid: true},
			Name:  sql.NullString{String: "user-name-1", Valid: true},
			Login: sql.NullString{String: "user-login-1", Valid: true},
			Role:  sql.NullString{String: "student", Valid: true},
		},
		{
			Id:    sql.NullString{String: "user-id-2", Valid: true},
			Name:  sql.NullString{String: "user-name-2", Valid: true},
			Login: sql.NullString{String: "user-login-2", Valid: true},
			Role:  sql.NullString{String: "teacher", Valid: true},
		},
	}

	// Вызов тестируемой функции
	entities := usersMapper.UsersInfoFromRepoToService(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Id.String, entities[i].Id)
		assert.Equal(t, repoModel.Name.String, entities[i].Name)
		assert.Equal(t, repoModel.Login.String, entities[i].Login)
		assert.Equal(t, usersModels.UserRole(repoModel.Role.String), entities[i].Role)
	}
}
