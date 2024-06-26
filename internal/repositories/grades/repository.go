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
	gradesTableQuery := `
		INSERT INTO grades (id, evaluation, subject_id, student_id, is_final
	`

	gradesTableQueryArgs := pgx.NamedArgs{
		"id":         grade.Id,
		"evaluation": grade.Evaluation,
		"subject_id": grade.SubjectId,
		"student_id": grade.StudentId,
		"is_final":   *grade.IsFinal,
	}

	if grade.Comment != nil {
		gradesTableQuery = fmt.Sprintf("%s, comment) VALUES (@id, @evaluation, @subject_id, @student_id, @is_final, @comment)", gradesTableQuery)
		gradesTableQueryArgs["comment"] = *grade.Comment
	} else {
		gradesTableQuery = fmt.Sprintf("%s ) VALUES (@id, @evaluation, @subject_id, @student_id, @is_final)", gradesTableQuery)
	}

	gradesTeachersQuery := `
		INSERT INTO grades_teachers (grade_id, teacher_id) 
		VALUES (@grade_id, @teacher_id)
	`

	gradesTeachersQueryArgs := pgx.NamedArgs{
		"grade_id":   grade.Id,
		"teacher_id": grade.TeacherId,
	}

	log.Printf("Creating grade: %#v", grade)

	tx, err := r.db.Begin(c.Context())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(c.Context()); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	if _, err = tx.Exec(c.Context(), gradesTableQuery, gradesTableQueryArgs); err != nil {
		log.Printf("Error creating grade: %#v, error: %v", grade, err)
		return err
	}

	if _, err = tx.Exec(c.Context(), gradesTeachersQuery, gradesTeachersQueryArgs); err != nil {
		log.Printf("Error add new record in grade_teachers error: %v", err)
		return err
	}

	if err = tx.Commit(c.Context()); err != nil {
		return err
	}

	log.Printf("Grade created successfully: %#v", grade)

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
	SELECT u.uuid, u.login, u.name, u.last_name, u.middle_name, g.id, g.evaluation, g.created_at, g.is_final, g.comment
	FROM groups gr
	JOIN users_groups ug ON ug.group_id = gr.id
	JOIN users u ON u.uuid = ug.user_id
	LEFT JOIN grades g ON g.student_id = u.uuid AND g.subject_id=$1
	`

	args := []interface{}{
		subjectId,
		groupId,
	}
	if isFinal != nil {
		query = fmt.Sprintf("%s AND g.is_final=$3", query)
		args = append(args, *isFinal)
	}
	query = fmt.Sprintf("%s WHERE gr.id = $2 ORDER BY uuid;", query)

	rows, err := r.db.Query(c.Context(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersSubjectGrades []gradesModels.UserSubjectGradesRepoModel

	for rows.Next() {
		var studentInfo usersModels.UserInfoRepoModel
		var gradeInfo gradesModels.GradeInfoRepoModel

		if err := rows.Scan(
			&studentInfo.Id,
			&studentInfo.Login,
			&studentInfo.Name,
			&studentInfo.LastName,
			&studentInfo.MiddleName,
			&gradeInfo.Id,
			&gradeInfo.Evaluation,
			&gradeInfo.CreatedAt,
			&gradeInfo.IsFinal,
			&gradeInfo.Comment,
		); err != nil {
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

func (r *repository) GetAllGradesInfo(c *fiber.Ctx) ([]gradesModels.GradesReportRecordRepoModel, error) {
	query := `
		SELECT u.last_name, u.name, u.middle_name, gr.number, s.name, g.evaluation, g.created_at, g.comment, g.is_final, t.last_name, t.name, t.middle_name
		FROM grades g
		JOIN users u ON g.student_id = u.uuid
		JOIN users_groups ug ON u.uuid = ug.user_id
		JOIN groups gr ON ug.group_id = gr.id
		JOIN subjects s ON g.subject_id = s.id
		JOIN grades_teachers gt ON g.id = gt.grade_id
		JOIN users t ON gt.teacher_id = t.uuid
		ORDER BY u.last_name, u.name, u.middle_name, s.name, g.created_at
	`

	rows, err := r.db.Query(c.Context(), query)
	if err != nil {
		log.Printf("Error fetching all grades: %v", err)
		return nil, err
	}
	defer rows.Close()

	var grades []gradesModels.GradesReportRecordRepoModel
	for rows.Next() {
		var grade gradesModels.GradesReportRecordRepoModel
		if err := rows.Scan(
			&grade.Student.LastName,
			&grade.Student.Name,
			&grade.Student.MiddleName,
			&grade.GroupNumber,
			&grade.GradeInfo.SubjectName,
			&grade.GradeInfo.Evaluation,
			&grade.GradeInfo.CreatedAt,
			&grade.GradeInfo.Comment,
			&grade.GradeInfo.IsFinal,
			&grade.Teacher.LastName,
			&grade.Teacher.Name,
			&grade.Teacher.MiddleName,
		); err != nil {
			log.Printf("Error scanning grade: %v", err)
			continue
		}
		log.Printf("Grade: %#v", grade)
		if grade.GradeInfo.Evaluation.Valid {
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
