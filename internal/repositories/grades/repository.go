package grades

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	gradesModels "sad/internal/models/grades"
	usersModels "sad/internal/models/users"

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
	query := `
		INSERT INTO grades (id, evaluation, subject_id, student_id, is_final, comment) 
		VALUES (@id, @evaluation, @subject_id, @student_id, @is_final, @comment)
	`
	log.Printf("Creating grade: %#v", grade)

	args := pgx.NamedArgs{
		"id":         grade.Id,
		"evaluation": grade.Evaluation,
		"subject_id": grade.SubjectId,
		"student_id": grade.StudentId,
		"is_final":   *grade.IsFinal,
		"comment":    *grade.Comment,
	}

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error creating grade: %#v, error: %v", grade, err)
	} else {
		log.Printf("Grade created successfully: %#v", grade)
	}

	return err
}

func (r *repository) GetAllStudentGrades(c *fiber.Ctx, studentId string, isFinal bool, subjectId *string) ([]gradesModels.GradeInfoRepoModel, error) {
	query := `
	SELECT g.id, s.name, evaluation, created_at, g.is_final, g.comment
	FROM grades g 
	JOIN subjects s on g.subject_id = s.id 
	WHERE student_id=$1 AND is_final=$2
	`
	args := []interface{}{
		studentId,
		isFinal,
	}
	if subjectId != nil {
		query = fmt.Sprintf("%s AND subject_id=$3", query)
		args = append(args, *subjectId)
	}
	query = fmt.Sprintf("%s ORDER BY created_at DESC", query)
	log.Printf("Fetching all student's %s grades", studentId)

	rows, err := r.db.Query(c.Context(), query, args...)
	if err != nil {
		log.Printf("Error fetching all grades: %v", err)
		return nil, err
	}
	defer rows.Close()

	var grades []gradesModels.GradeInfoRepoModel
	for rows.Next() {
		var grade gradesModels.GradeInfoRepoModel
		if err := rows.Scan(
			&grade.Id,
			&grade.SubjectName,
			&grade.Evaluation,
			&grade.CreatedAt,
			&grade.IsFinal,
			&grade.Comment); err != nil {
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

	log.Printf("All grades fetched successfully %#v", grades)
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
	query := `
		UPDATE grades 
		SET evaluation=@evaluation, comment=@comment 
		WHERE id=@grade_id
	`

	args := pgx.NamedArgs{
		"evaluation": grade.Evaluation,
		"grade_id":   grade.Id,
		"comment":    *grade.Comment,
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
	query := `
		SELECT id, subject_id, student_id, evaluation, created_at, comment 
		FROM grades 
		WHERE id=$1
	`
	log.Printf("Fetching grade by id: %s", gradeId)

	row := r.db.QueryRow(c.Context(), query, gradeId)

	grade := &gradesModels.GradeRepoModel{}
	err := row.Scan(
		&grade.Id,
		&grade.SubjectId,
		&grade.StudentId,
		&grade.Evaluation,
		&grade.CreatedAt,
		&grade.Comment,
	)
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

func (r *repository) GetStudentsGradesBySubjectAndGroup(c *fiber.Ctx, subjectId, groupId string, isFinal *bool) ([]gradesModels.UserSubjectGradesRepoModel, error) {
	query := `
	SELECT u.uuid, u.login, u.name, g.id, g.evaluation, g.created_at, g.is_final, g.comment
	FROM groups gr
	JOIN users_groups ug ON ug.group_id = gr.id
	JOIN users u ON u.uuid = ug.user_id
	LEFT JOIN grades g ON g.student_id = u.uuid AND g.subject_id=$2 AND g.is_final=$3
	WHERE gr.id = $1
	`

	args := []interface{}{
		subjectId,
		groupId,
	}
	if isFinal != nil {
		query = fmt.Sprintf("%s AND is_final=$3", query)
		args = append(args, *isFinal)
	}
	query = fmt.Sprintf("%s ORDER BY uuid", query)

	rows, err := r.db.Query(c.Context(), query, groupId, subjectId, isFinal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersSubjectGrades []gradesModels.UserSubjectGradesRepoModel

	for rows.Next() {
		var studentInfo usersModels.UserInfoRepoModel
		var gradeInfo gradesModels.GradeInfoRepoModel

		err := rows.Scan(
			&studentInfo.Id,
			&studentInfo.Login,
			&studentInfo.Name,
			&gradeInfo.Id,
			&gradeInfo.Evaluation,
			&gradeInfo.CreatedAt,
			&gradeInfo.IsFinal,
			&gradeInfo.Comment,
		)
		if err != nil {
			return nil, err
		}

		var userExist bool
		if gradeInfo.Id.Valid {
			for i := range usersSubjectGrades {
				if usersSubjectGrades[i].Student.Id.String == studentInfo.Id.String {
					userExist = true
					usersSubjectGrades[i].Grades = append(usersSubjectGrades[i].Grades, gradeInfo)
					break
				}
			}
		}

		if !userExist {
			grades := make([]gradesModels.GradeInfoRepoModel, 0)
			if gradeInfo.Id.Valid {
				grades = append(grades, gradeInfo)
			}
			usersSubjectGrades = append(usersSubjectGrades, gradesModels.UserSubjectGradesRepoModel{
				Student: studentInfo,
				Grades:  grades,
			})
		}

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Grades %#v", usersSubjectGrades)
	return usersSubjectGrades, nil
}
