package groups

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	groupsMapper "sad/internal/mappers/groups"
	groupsModels "sad/internal/models/groups"
	subjectsModels "sad/internal/models/subjects"
	usersModels "sad/internal/models/users"
	"testing"
)

func TestFromGroupDetailsRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := groupsModels.GroupDetailsRepo{
		GroupRepoModel: groupsModels.GroupRepoModel{
			Id:     sql.NullString{String: "group-id", Valid: true},
			Number: sql.NullString{String: "group-number", Valid: true},
		},
		Users: []usersModels.UserInfoRepoModel{
			{
				Id:    sql.NullString{String: "user-id-1", Valid: true},
				Login: sql.NullString{String: "user-login-1", Valid: true},
				Name:  sql.NullString{String: "user-first-name-1", Valid: true},
				Role:  sql.NullString{String: "student", Valid: true},
			},
			{
				Id:    sql.NullString{String: "user-id-2", Valid: true},
				Login: sql.NullString{String: "user-login-2", Valid: true},
				Name:  sql.NullString{String: "user-first-name-2", Valid: true},
				Role:  sql.NullString{String: "student", Valid: true},
			},
		},
	}

	// Вызов тестируемой функции
	entity := groupsMapper.FromGroupDetailsRepoModelToEntity(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Group.Id)
	assert.Equal(t, repoModel.Number.String, entity.Group.Number)
	assert.Equal(t, len(repoModel.Users), len(entity.Users))
	for i, userRepo := range repoModel.Users {
		assert.Equal(t, userRepo.Id.String, entity.Users[i].Id)
		assert.Equal(t, userRepo.Login.String, entity.Users[i].Login)
		assert.Equal(t, userRepo.Name.String, entity.Users[i].Name)
		assert.Equal(t, usersModels.UserRole(userRepo.Role.String), entity.Users[i].Role)
	}
}

func TestFromGroupRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := groupsModels.GroupRepoModel{
		Id:     sql.NullString{String: "group-id", Valid: true},
		Number: sql.NullString{String: "group-number", Valid: true},
	}

	// Вызов тестируемой функции
	entity := groupsMapper.FromGroupRepoModelToEntity(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.Number.String, entity.Number)
}

func TestFromGroupsRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []groupsModels.GroupRepoModel{
		{
			Id:     sql.NullString{String: "group-id-1", Valid: true},
			Number: sql.NullString{String: "group-number-1", Valid: true},
		},
		{
			Id:     sql.NullString{String: "group-id-2", Valid: true},
			Number: sql.NullString{String: "group-number-2", Valid: true},
		},
	}

	// Вызов тестируемой функции
	entities := groupsMapper.FromGroupsRepoModelToEntity(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Id.String, entities[i].Id)
		assert.Equal(t, repoModel.Number.String, entities[i].Number)
	}
}

func TestGroupsWithSubjectsFromRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []subjectsModels.GroupsWithSubjectsRepoModel{
		{
			Group: groupsModels.GroupRepoModel{
				Id:     sql.NullString{String: "group-id-1", Valid: true},
				Number: sql.NullString{String: "group-number-1", Valid: true},
			},
			Subjects: []subjectsModels.SubjectRepoModel{
				{
					Id:   sql.NullString{String: "subject-id-1", Valid: true},
					Name: sql.NullString{String: "subject-name-1", Valid: true},
				},
				{
					Id:   sql.NullString{String: "subject-id-2", Valid: true},
					Name: sql.NullString{String: "subject-name-2", Valid: true},
				},
			},
		},
		{
			Group: groupsModels.GroupRepoModel{
				Id:     sql.NullString{String: "group-id-2", Valid: true},
				Number: sql.NullString{String: "group-number-2", Valid: true},
			},
			Subjects: []subjectsModels.SubjectRepoModel{
				{
					Id:   sql.NullString{String: "subject-id-3", Valid: true},
					Name: sql.NullString{String: "subject-name-3", Valid: true},
				},
				{
					Id:   sql.NullString{String: "subject-id-4", Valid: true},
					Name: sql.NullString{String: "subject-name-4", Valid: true},
				},
			},
		},
	}

	// Вызов тестируемой функции
	entities := groupsMapper.GroupsWithSubjectsFromRepoModelToEntity(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Group.Id.String, entities[i].Group.Id)
		assert.Equal(t, repoModel.Group.Number.String, entities[i].Group.Number)
		assert.Equal(t, len(repoModel.Subjects), len(entities[i].Subjects))
		for j, subjectRepo := range repoModel.Subjects {
			assert.Equal(t, subjectRepo.Id.String, entities[i].Subjects[j].Id)
			assert.Equal(t, subjectRepo.Name.String, entities[i].Subjects[j].Name)
		}
	}
}
