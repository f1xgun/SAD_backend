package subjects

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	subjectsMappers "sad/internal/mappers/subjects"
	subjectsModels "sad/internal/models/subjects"
	usersModels "sad/internal/models/users"
	"testing"
)

func TestFromSubjectRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := subjectsModels.SubjectRepoModel{
		Id:   sql.NullString{String: "subject-id", Valid: true},
		Name: sql.NullString{String: "subject-name", Valid: true},
	}

	// Вызов тестируемой функции
	entity := subjectsMappers.FromSubjectRepoModelToEntity(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.Name.String, entity.Name)
}

func TestFromSubjectsRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []subjectsModels.SubjectRepoModel{
		{
			Id:   sql.NullString{String: "subject-id-1", Valid: true},
			Name: sql.NullString{String: "subject-name-1", Valid: true},
		},
		{
			Id:   sql.NullString{String: "subject-id-2", Valid: true},
			Name: sql.NullString{String: "subject-name-2", Valid: true},
		},
	}

	// Вызов тестируемой функции
	entities := subjectsMappers.FromSubjectsRepoModelToEntity(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Id.String, entities[i].Id)
		assert.Equal(t, repoModel.Name.String, entities[i].Name)
	}
}

func TestFromSubjectDetailsRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := subjectsModels.SubjectInfoRepoModel{
		Id:   sql.NullString{String: "subject-id", Valid: true},
		Name: sql.NullString{String: "subject-name", Valid: true},
		Teacher: usersModels.UserInfoRepoModel{
			Id:    sql.NullString{String: "teacher-id", Valid: true},
			Login: sql.NullString{String: "teacher-login", Valid: true},
			Name:  sql.NullString{String: "teacher-first-name", Valid: true},
			Role:  sql.NullString{String: "teacher", Valid: true},
		},
	}

	// Вызов тестируемой функции
	entity := subjectsMappers.FromSubjectDetailsRepoModelToEntity(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.Name.String, entity.Name)
	assert.Equal(t, repoModel.Teacher.Id.String, entity.Teacher.Id)
	assert.Equal(t, repoModel.Teacher.Login.String, entity.Teacher.Login)
	assert.Equal(t, repoModel.Teacher.Name.String, entity.Teacher.Name)
	assert.Equal(t, usersModels.UserRole(repoModel.Teacher.Role.String), entity.Teacher.Role)
}
