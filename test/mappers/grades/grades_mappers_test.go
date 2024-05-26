package grades

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sad/internal/mappers/grades"
	gradesModels "sad/internal/models/grades"
	usersModels "sad/internal/models/users"
	"testing"
	"time"
)

func TestFromGradeRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := gradesModels.GradeRepoModel{
		Id:         sql.NullString{String: uuid.New().String(), Valid: true},
		StudentId:  sql.NullString{String: uuid.New().String(), Valid: true},
		SubjectId:  sql.NullString{String: uuid.New().String(), Valid: true},
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		Evaluation: sql.NullInt16{Int16: 5, Valid: true},
		TeacherId:  sql.NullString{String: uuid.New().String(), Valid: true},
	}

	// Вызов тестируемой функции
	entity := grades.FromGradeRepoModelToEntity(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.StudentId.String, entity.StudentId)
	assert.Equal(t, repoModel.SubjectId.String, entity.SubjectId)
	assert.Equal(t, repoModel.CreatedAt.Time, entity.CreatedAt)
	assert.Equal(t, int(repoModel.Evaluation.Int16), entity.Evaluation)
	assert.Equal(t, repoModel.TeacherId.String, entity.TeacherId)
}

func TestFromGradesRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []gradesModels.GradeRepoModel{
		{
			Id:         sql.NullString{String: uuid.New().String(), Valid: true},
			StudentId:  sql.NullString{String: uuid.New().String(), Valid: true},
			SubjectId:  sql.NullString{String: uuid.New().String(), Valid: true},
			CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			Evaluation: sql.NullInt16{Int16: 5, Valid: true},
			TeacherId:  sql.NullString{String: uuid.New().String(), Valid: true},
		},
		{
			Id:         sql.NullString{String: uuid.New().String(), Valid: true},
			StudentId:  sql.NullString{String: uuid.New().String(), Valid: true},
			SubjectId:  sql.NullString{String: uuid.New().String(), Valid: true},
			CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			Evaluation: sql.NullInt16{Int16: 4, Valid: true},
			TeacherId:  sql.NullString{String: uuid.New().String(), Valid: true},
		},
	}

	// Вызов тестируемой функции
	entities := grades.FromGradesRepoModelToEntity(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Id.String, entities[i].Id)
		assert.Equal(t, repoModel.StudentId.String, entities[i].StudentId)
		assert.Equal(t, repoModel.SubjectId.String, entities[i].SubjectId)
		assert.Equal(t, repoModel.CreatedAt.Time, entities[i].CreatedAt)
		assert.Equal(t, int(repoModel.Evaluation.Int16), entities[i].Evaluation)
		assert.Equal(t, repoModel.TeacherId.String, entities[i].TeacherId)
	}
}

func TestFromGradeInfoRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModel := gradesModels.GradeInfoRepoModel{
		Id:          sql.NullString{String: uuid.New().String(), Valid: true},
		SubjectName: sql.NullString{String: "Math", Valid: true},
		TeacherName: sql.NullString{String: "John Doe", Valid: true},
		Evaluation:  sql.NullInt16{Int16: 5, Valid: true},
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		IsFinal:     sql.NullBool{Bool: true, Valid: true},
	}

	// Вызов тестируемой функции
	entity := grades.FromGradeInfoRepoModelToEntity(repoModel)

	// Проверка результатов
	assert.Equal(t, repoModel.Id.String, entity.Id)
	assert.Equal(t, repoModel.SubjectName.String, entity.SubjectName)
	assert.Equal(t, repoModel.TeacherName.String, entity.TeacherName)
	assert.Equal(t, int(repoModel.Evaluation.Int16), entity.Evaluation)
	assert.Equal(t, repoModel.CreatedAt.Time, entity.CreatedAt)
	assert.NotNil(t, entity.IsFinal)
	assert.Equal(t, repoModel.IsFinal.Bool, *entity.IsFinal)
}

func TestFromGradesInfoRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []gradesModels.GradeInfoRepoModel{
		{
			Id:          sql.NullString{String: uuid.New().String(), Valid: true},
			SubjectName: sql.NullString{String: "Math", Valid: true},
			TeacherName: sql.NullString{String: "John Doe", Valid: true},
			Evaluation:  sql.NullInt16{Int16: 5, Valid: true},
			CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			IsFinal:     sql.NullBool{Bool: true, Valid: true},
		},
		{
			Id:          sql.NullString{String: uuid.New().String(), Valid: true},
			SubjectName: sql.NullString{String: "English", Valid: true},
			TeacherName: sql.NullString{String: "Jane Smith", Valid: true},
			Evaluation:  sql.NullInt16{Int16: 4},
			CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			IsFinal:     sql.NullBool{Bool: false, Valid: true},
		},
	}

	// Вызов тестируемой функции
	entities := grades.FromGradesInfoRepoModelToEntity(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Id.String, entities[i].Id)
		assert.Equal(t, repoModel.SubjectName.String, entities[i].SubjectName)
		assert.Equal(t, repoModel.TeacherName.String, entities[i].TeacherName)
		assert.Equal(t, int(repoModel.Evaluation.Int16), entities[i].Evaluation)
		assert.Equal(t, repoModel.CreatedAt.Time, entities[i].CreatedAt)
		assert.NotNil(t, entities[i].IsFinal)
		assert.Equal(t, repoModel.IsFinal.Bool, *entities[i].IsFinal)
	}
}

func TestFromUserWithGradesRepoModelToEntity(t *testing.T) {
	// Подготовка тестовых данных
	repoModels := []gradesModels.UserSubjectGradesRepoModel{
		{
			Student: usersModels.UserInfoRepoModel{
				Id:    sql.NullString{String: uuid.New().String(), Valid: true},
				Login: sql.NullString{String: "student1", Valid: true},
				Name:  sql.NullString{String: "John", Valid: true},
				Role:  sql.NullString{String: "student", Valid: true},
			},
			Grades: []gradesModels.GradeInfoRepoModel{
				{
					Id:          sql.NullString{String: uuid.New().String(), Valid: true},
					SubjectName: sql.NullString{String: "Math", Valid: true},
					TeacherName: sql.NullString{String: "John Doe", Valid: true},
					Evaluation:  sql.NullInt16{Int16: 5, Valid: true},
					CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
					IsFinal:     sql.NullBool{Bool: true, Valid: true},
				},
			},
		},
		{
			Student: usersModels.UserInfoRepoModel{
				Id:    sql.NullString{String: uuid.New().String(), Valid: true},
				Login: sql.NullString{String: "student2", Valid: true},
				Name:  sql.NullString{String: "Jane", Valid: true},
				Role:  sql.NullString{String: "student", Valid: true},
			},
			Grades: []gradesModels.GradeInfoRepoModel{
				{
					Id:          sql.NullString{String: uuid.New().String(), Valid: true},
					SubjectName: sql.NullString{String: "English", Valid: true},
					TeacherName: sql.NullString{String: "Jane Smith", Valid: true},
					Evaluation:  sql.NullInt16{Int16: 4},
					CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
					IsFinal:     sql.NullBool{Bool: false, Valid: true},
				},
			},
		},
	}

	// Вызов тестируемой функции
	entities := grades.FromUserWithGradesRepoModelToEntity(repoModels)

	// Проверка результатов
	assert.Equal(t, len(repoModels), len(entities))
	for i, repoModel := range repoModels {
		assert.Equal(t, repoModel.Student.Id.String, entities[i].Student.Id)
		assert.Equal(t, repoModel.Student.Login.String, entities[i].Student.Login)
		assert.Equal(t, repoModel.Student.Name.String, entities[i].Student.Name)
		assert.Equal(t, usersModels.UserRole(repoModel.Student.Role.String), entities[i].Student.Role)

		assert.Equal(t, len(repoModel.Grades), len(entities[i].Grades))
		for j, gradeRepoModel := range repoModel.Grades {
			assert.Equal(t, gradeRepoModel.Id.String, entities[i].Grades[j].Id)
			assert.Equal(t, gradeRepoModel.SubjectName.String, entities[i].Grades[j].SubjectName)
			assert.Equal(t, gradeRepoModel.TeacherName.String, entities[i].Grades[j].TeacherName)
			assert.Equal(t, int(gradeRepoModel.Evaluation.Int16), entities[i].Grades[j].Evaluation)
			assert.Equal(t, gradeRepoModel.CreatedAt.Time, entities[i].Grades[j].CreatedAt)
			assert.NotNil(t, entities[i].Grades[j].IsFinal)
			assert.Equal(t, gradeRepoModel.IsFinal.Bool, *entities[i].Grades[j].IsFinal)
		}
	}
}
