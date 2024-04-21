package grades

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	gradesModels "sad/internal/models/grades"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(c *fiber.Ctx, grade gradesModels.Grade) error {
	query := "INSERT INTO grades (id, evaluation, subject_id, student_id) VALUES (@id, @evaluation, @subject_id, @student_id)"
	log.Printf("Creating grade: %#v", grade)

	args := pgx.NamedArgs{
		"id":         grade.Id,
		"evaluation": grade.Evaluation,
		"subject_id": grade.SubjectId,
		"student_id": grade.StudentId,
	}
	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error creating grade: %#v, error: %v", grade, err)
	} else {
		log.Printf("Grade created successfully: %#v", grade)
	}

	return err
}

func (r *repository) GetAllStudentGrades(c *fiber.Ctx, studentId string) ([]gradesModels.GradeInfoRepoModel, error) {
	query := "SELECT g.id, s.name, evaluation, created_at FROM grades g JOIN subjects s on g.subject_id = s.id WHERE student_id=$1 ORDER BY created_at DESC"
	log.Printf("Fetching all student's %s grades", studentId)

	rows, err := r.db.Query(c.Context(), query, studentId)
	if err != nil {
		log.Printf("Error fetching all grades: %v", err)
		return nil, err
	}
	defer rows.Close()

	var grades []gradesModels.GradeInfoRepoModel
	for rows.Next() {
		var grade gradesModels.GradeInfoRepoModel
		if err := rows.Scan(&grade.Id, &grade.SubjectName, &grade.Evaluation, &grade.CreatedAt); err != nil {
			log.Printf("Error scanning grade: %v", err)
			continue
		}
		if grade.Id.Valid {
			grades = append(grades, grade)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over grades: %v", err)
		return nil, err
	}

	log.Printf("All grades fetched successfully")
	return grades, nil
}

func (r *repository) Delete(c *fiber.Ctx, gradeId string) error {
	query := "DELETE FROM grades WHERE id=$1"
	log.Printf("Deleting grade %s", gradeId)

	_, err := r.db.Exec(c.Context(), query, gradeId)
	if err != nil {
		log.Printf("Error deleting grade error: %v", err)
		return err
	}

	log.Printf("Grade deleted successfully: grade_id=%s", gradeId)
	return nil
}

func (r *repository) Update(c *fiber.Ctx, grade gradesModels.Grade) error {
	query := "UPDATE grades SET evaluation=@evaluation WHERE id=@grade_id"
	args := pgx.NamedArgs{
		"evaluation": grade.Evaluation,
		"grade_id":   grade.Id,
	}

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error updating grade with id %s, err: %v", grade.Id, err)
	} else {
		log.Printf("Update grade with id %s successfully", grade.Id)
	}

	return err
}

func (r *repository) GetById(c *fiber.Ctx, gradeId string) (*gradesModels.GradeRepoModel, error) {
	query := "SELECT id, subject_id, student_id, evaluation, created_at FROM grades WHERE id=$1"
	log.Printf("Fetching grade by id: %s", gradeId)

	row := r.db.QueryRow(c.Context(), query, gradeId)

	grade := &gradesModels.GradeRepoModel{}
	err := row.Scan(&grade.Id, &grade.SubjectId, &grade.StudentId, &grade.Evaluation, &grade.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Grade not found with id: %s", gradeId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching grade by id: %s, error: %v", gradeId, err)
		return nil, err
	}

	log.Printf("Grade fetched successfully by id: %s", gradeId)
	return grade, nil
}
