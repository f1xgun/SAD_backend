package grades

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"log"
	gradesMapper "sad/internal/mappers/grades"
	errorsModels "sad/internal/models/errors"
	gradesModels "sad/internal/models/grades"
	"sad/internal/repositories"

	"github.com/google/uuid"
)

type service struct {
	gradesRepository repositories.GradesRepository
	usersRepository  repositories.UserRepository
}

func NewService(gradesRepository repositories.GradesRepository, usersRepository repositories.UserRepository) *service {
	return &service{
		gradesRepository: gradesRepository,
		usersRepository:  usersRepository,
	}
}

func (s *service) Create(ctx context.Context, grade gradesModels.Grade) error {
	log.Printf("Creating a new grade: %#v", grade)

	if grade.TeacherId == "" {
		return errors.New("teacher_id is required")
	}

	if grade.StudentId == "" {
		return errors.New("student_id is required")
	}

	if grade.SubjectId == "" {
		return errors.New("subject_id is required")
	}

	newGrade := gradesModels.Grade{
		Id:         uuid.New().String(),
		StudentId:  grade.StudentId,
		SubjectId:  grade.SubjectId,
		Evaluation: grade.Evaluation,
		IsFinal:    grade.IsFinal,
		Comment:    grade.Comment,
		TeacherId:  grade.TeacherId,
	}

	if err := s.gradesRepository.Create(ctx, newGrade); err != nil {
		log.Printf("Error creating grade '%s' in the repository: %s", newGrade.Id, err.Error())
		return errorsModels.ErrServer
	}

	return nil
}

func (s *service) Delete(ctx context.Context, gradeId string) error {
	log.Printf("Attempting to delete grade '%s'.", gradeId)

	grade, err := s.gradesRepository.GetById(ctx, gradeId)
	if err != nil {
		log.Printf("Error retrieving grade '%s': %v", gradeId, err)
		return err
	}

	if grade == nil {
		log.Printf("Grade '%s' does not exist.", gradeId)
		return errorsModels.ErrGradeDoesNotExist
	}

	if err := s.gradesRepository.Delete(ctx, gradeId); err != nil {
		log.Printf("Error deleting grade '%s': %v", gradeId, err)
		return err
	}

	log.Printf("Grade '%s' successfully deleted.", gradeId)
	return nil
}

func (s *service) Update(ctx context.Context, gradeId string, evaluation *int, comment *string) error {
	log.Printf("Attempting to update grade '%s'.", gradeId)

	existedGrade, err := s.gradesRepository.GetById(ctx, gradeId)
	if err != nil {
		log.Printf("Error retrieving grade '%s': %v", gradeId, err)
		return err
	}

	if existedGrade == nil || !existedGrade.Id.Valid {
		log.Printf("Grade '%s' does not exist.", gradeId)
		return errorsModels.ErrGradeDoesNotExist
	}

	updatedExistedGrade := gradesMapper.FromGradeRepoModelToEntity(*existedGrade)

	if evaluation != nil {
		updatedExistedGrade.Evaluation = *evaluation
	}

	if comment != nil {
		updatedExistedGrade.Comment = comment
	}

	if err := s.gradesRepository.Update(ctx, updatedExistedGrade); err != nil {
		log.Printf("Error updating grade '%s': %v", gradeId, err)
		return err
	}

	log.Printf("Grade '%s' successfully updated.", gradeId)
	return nil
}

func (s *service) GetAllStudentGrades(ctx context.Context, studentId string, isFinal bool, subjectId *string) ([]gradesModels.GradeInfo, error) {
	log.Printf("Attemting to get student's with id '%s' grades", studentId)

	isUserExist, err := s.usersRepository.CheckUserExists(ctx, studentId)
	if err != nil {
		log.Printf("Error retrieving student with id '%s'", studentId)
		return nil, errorsModels.ErrServer
	}

	if !isUserExist {
		log.Printf("Error student with id '%s' doesn't exist", studentId)
		return nil, errorsModels.ErrUserDoesNotExist
	}

	gradesRepo, err := s.gradesRepository.GetAllStudentGrades(ctx, studentId, isFinal, subjectId)

	if err != nil {
		log.Printf("Error retrieving grades")
		return nil, err
	}

	grades := gradesMapper.FromGradesInfoRepoModelToEntity(gradesRepo)

	return grades, nil
}

func (s *service) GetStudentsGradesBySubjectAndGroup(ctx context.Context, subjectId, groupId string, isFinal *bool) ([]gradesModels.UserSubjectGrades, error) {
	userWithGradesRepo, err := s.gradesRepository.GetStudentsGradesBySubjectAndGroup(ctx, subjectId, groupId, isFinal)

	if err != nil {
		return nil, err
	}

	userWithGrades := gradesMapper.FromUserWithGradesRepoModelToEntity(userWithGradesRepo)

	return userWithGrades, nil
}

func (s *service) GetGradesInCsv(ctx context.Context) (string, error) {
	gradesReportRecords, err := s.gradesRepository.GetAllGradesInfo(ctx)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	columnNames := []string{
		"Фамилия студента",
		"Имя студента",
		"Отчество студента",
		"Номер группы",
		"Название предмета",
		"Оценка",
		"Дата выставления оценки",
		"Комментарий",
		"Итоговая ли оценка",
		"Фамилия преподавателя",
		"Имя преподавателя",
		"Отчество преподавателя",
	}

	if err := writer.Write(columnNames); err != nil {
		log.Printf("Error writing record to csv: %v", err)
		return "", err
	}

	for _, grade := range gradesReportRecords {
		if err := writer.Write(grade.ToCsvString()); err != nil {
			log.Printf("Error writing record to csv: %#v", err)
			return "", err
		}
	}

	if err := writer.Error(); err != nil {
		return "", err
	}

	writer.Flush()

	return buffer.String(), nil
}
