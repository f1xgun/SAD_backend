package subjects

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	subjectsModels "sad/internal/models/subjects"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(c *fiber.Ctx, subject subjectsModels.Subject) error {
	query := "INSERT INTO subjects (id, name) VALUES (@id, @name)"
	log.Printf("Creating subject: %#v", subject)

	args := pgx.NamedArgs{
		"id":   subject.Id,
		"name": subject.Name,
	}
	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error creating subject: %#v, error: %v", subject, err)
	} else {
		log.Printf("Subject created successfully: %#v", subject)
	}

	return err
}

func (r *repository) GetAll(c *fiber.Ctx) ([]subjectsModels.SubjectRepoModel, error) {
	query := "SELECT id, name FROM subjects"
	log.Printf("Fetching all subjects")

	rows, err := r.db.Query(c.Context(), query)
	if err != nil {
		log.Printf("Error fetching all subjects: %v", err)
		return nil, err
	}
	defer rows.Close()

	var subjects []subjectsModels.SubjectRepoModel
	for rows.Next() {
		var subject subjectsModels.SubjectRepoModel
		if err := rows.Scan(&subject.Id, &subject.Name); err != nil {
			log.Printf("Error scanning subject: %v", err)
			continue
		}

		if subject.Id.Valid {
			subjects = append(subjects, subject)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over subjects: %v", err)
		return nil, err
	}

	log.Printf("All subjects fetched successfully")
	return subjects, nil
}

func (r *repository) GetById(c *fiber.Ctx, subjectId string) (*subjectsModels.SubjectRepoModel, error) {
	query := "SELECT id, name FROM subjects WHERE id=$1"
	log.Printf("Fetching subject by id: %s", subjectId)

	row := r.db.QueryRow(c.Context(), query, subjectId)

	subject := &subjectsModels.SubjectRepoModel{}
	err := row.Scan(&subject.Id, &subject.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Subject not found with id: %s", subjectId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching subject by id: %s, error: %v", subjectId, err)
		return nil, err
	}

	log.Printf("Subject fetched successfully by id: %s", subjectId)
	return subject, nil
}

func (r *repository) AddSubjectToGroup(c *fiber.Ctx, subjectGroup subjectsModels.SubjectGroup) error {
	log.Printf("Adding subject %s to group %s", subjectGroup.SubjectId, subjectGroup.GroupId)
	subjectTeacherId, err := r.GetSubjectTeacherId(c, subjectGroup.SubjectId, subjectGroup.TeacherId)
	if err != nil {
		log.Printf("Error adding subject to group: subject_id=%s, group_id=%s, error: %v",
			subjectGroup.SubjectId, subjectGroup.GroupId, err)
		return err
	}

	if subjectTeacherId == "" {
		log.Printf("Record with subject id %s and teacher id %s not found",
			subjectGroup.SubjectId, subjectGroup.TeacherId)
	}

	query := "INSERT INTO groups_subjects (group_id, subject_teacher_id) VALUES (@group_id, @subject_teacher_id)"
	args := pgx.NamedArgs{
		"subject_teacher_id": subjectTeacherId,
		"group_id":           subjectGroup.GroupId,
	}

	_, err = r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error adding subject to group: subject_id=%s, group_id=%s, teacher_id=%s, error: %v",
			subjectGroup.SubjectId, subjectGroup.GroupId, subjectGroup.TeacherId, err)
		return err
	}

	log.Printf("Subject added to group successfully: subject_id=%s, group_id=%s, teacher_id=%s",
		subjectGroup.SubjectId, subjectGroup.GroupId, subjectGroup.TeacherId)
	return nil
}

func (r *repository) DeleteSubjectFromGroup(c *fiber.Ctx, subjectId, groupId string) error {
	query := "DELETE FROM groups_subjects WHERE subject_id=@subject_id AND group_id=@group_id"
	args := pgx.NamedArgs{
		"subject_id": subjectId,
		"group_id":   groupId,
	}
	log.Printf("Deleting subject %s from group %s", subjectId, groupId)

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error deleting subject from group: subject_id=%s, group_id=%s, error: %v", subjectId, groupId, err)
		return err
	}

	log.Printf("Subject deleted from group successfully: subject_id=%s, group_id=%s", subjectId, groupId)
	return nil
}

func (r *repository) DeleteSubject(c *fiber.Ctx, subjectId string) error {
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

	args := pgx.NamedArgs{
		"subject_id": subjectId,
	}
	query := "DELETE FROM subjects_teachers WHERE subject_id=@subject_id"

	if _, err := tx.Exec(c.Context(), query, args); err != nil {
		log.Printf("Error deleting groups from subject with id %s, err: %v", subjectId, err)
		return err
	}
	log.Printf("Delete groups from subject with id %s successfully", subjectId)

	query = "DELETE FROM subjects WHERE id=@subject_id"
	if _, err = r.db.Exec(c.Context(), query, args); err != nil {
		log.Printf("Error deleting subject with id %s , err: %v", subjectId, err)
		return err
	}
	log.Printf("Delete subject with id %s successfully", subjectId)

	if err = tx.Commit(c.Context()); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateSubject(c *fiber.Ctx, subject subjectsModels.Subject) error {
	log.Printf("Update subject with new value: %#v", subject)

	query := "UPDATE subjects SET name=@name WHERE id=@subjectId"
	args := pgx.NamedArgs{
		"subjectId": subject.Id,
		"name":      subject.Name,
	}

	if _, err := r.db.Exec(c.Context(), query, args); err != nil {
		log.Printf("Error updating subject with id %s, err: %v", subject.Id, err)
		return err
	}

	log.Printf("Update subject with id %s successfully", subject.Id)

	return nil
}

func (r *repository) IsSubjectInGroup(c *fiber.Ctx, subjectId, groupId string) (bool, error) {
	log.Printf("Checking if subject with ID '%s' is in group with ID '%s'\n", subjectId, groupId)

	query := "SELECT COUNT(*) FROM groups_subjects WHERE subject_id=@subject_id and group_id=@group_id"
	args := pgx.NamedArgs{
		"subject_id": subjectId,
		"group_id":   groupId,
	}

	var count int
	err := r.db.QueryRow(c.Context(), query, args).Scan(&count)
	if err != nil {
		log.Printf("Error checking if subject '%s' is in group '%s': %v\n", subjectId, groupId, err)
		return false, err
	}

	if count > 0 {
		log.Printf("Subject '%s' is in group '%s'\n", subjectId, groupId)
	} else {
		log.Printf("Subject '%s' is not in group '%s'\n", subjectId, groupId)
	}

	return count > 0, nil
}

func (r *repository) GetSubjectTeacherId(c *fiber.Ctx, subjectId, teacherId string) (string, error) {
	query := "SELECT id FROM subjects_teachers WHERE subject_id=@subject_id AND teacher_id=@teacher_id"
	log.Printf("Fetching subject by id: %s", subjectId)

	args := pgx.NamedArgs{
		"subject_id": subjectId,
		"teacher_id": teacherId,
	}

	row := r.db.QueryRow(c.Context(), query, args)

	var subjectTeacherId string
	err := row.Scan(&subjectTeacherId)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Record with subject id %s and teacher id %s not found", subjectId, teacherId)
		return subjectTeacherId, nil
	} else if err != nil {
		log.Printf("Error fetching record with subject id %s and teacher id %s, error: %v", subjectId, teacherId, err)
		return subjectTeacherId, err
	}

	log.Printf("Record with subject id %s and teacher id %s fetched successfully", subjectId, teacherId)
	return subjectTeacherId, nil
}

func (r *repository) AddTeacherToSubject(c *fiber.Ctx, subjectTeacherId, subjectId, teacherId string) error {
	log.Printf("Adding teacher %s to subject %s", teacherId, subjectId)

	queryForInsert := `
	INSERT INTO subjects_teachers (id, subject_id, teacher_id) 
	VALUES (@subjectTeacherId, @subjectId, @teacherId)
	`
	args := pgx.NamedArgs{
		"subjectTeacherId": subjectTeacherId,
		"subjectId":        subjectId,
		"teacherId":        teacherId,
	}

	if _, err := r.db.Exec(c.Context(), queryForInsert, args); err != nil {
		log.Printf("Error adding teacher to subject: subject_id=%s, teacher_id=%s, error: %v",
			subjectId, teacherId, err)
		return err
	}

	log.Printf("Teacher added to subject successfully: subject_id=%s, teacher_id=%s",
		subjectId, teacherId)
	return nil
}

func (r *repository) GetByIdWithDetails(c *fiber.Ctx, subjectId string) (*subjectsModels.SubjectInfoRepoModel, error) {
	query := `
		SELECT s.id, s.name, u.uuid, u.login, u.name
		FROM subjects s
		LEFT JOIN subjects_teachers st ON s.id = st.subject_id
		LEFT JOIN users u ON u.uuid = st.teacher_id
		WHERE s.id = $1;
	`

	row := r.db.QueryRow(c.Context(), query, subjectId)
	var subject subjectsModels.SubjectInfoRepoModel
	err := row.Scan(&subject.Id, &subject.Name, &subject.Teacher.Id, &subject.Teacher.Login, &subject.Teacher.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Subject with id %s not found", subjectId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching subject with details: %v", err)
		return nil, err
	}

	log.Printf("Subject with details fetched successfully")
	return &subject, nil
}
