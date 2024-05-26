package subjects

import (
	"errors"
	"log"
	subjectsMappers "sad/internal/mappers/subjects"
	usersMapper "sad/internal/mappers/users"
	usersModels "sad/internal/models/users"
	"sad/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	errorsModels "sad/internal/models/errors"
	subjectsModels "sad/internal/models/subjects"
)

type service struct {
	groupsRepository   repositories.GroupsRepository
	subjectsRepository repositories.SubjectsRepository
	usersRepository    repositories.UserRepository
}

func NewService(
	groupsRepository repositories.GroupsRepository,
	subjectsRepository repositories.SubjectsRepository,
	usersRepository repositories.UserRepository,
) *service {
	return &service{
		groupsRepository:   groupsRepository,
		subjectsRepository: subjectsRepository,
		usersRepository:    usersRepository,
	}
}

func (s *service) Create(c *fiber.Ctx, name string, teacherId string) error {
	log.Println("Creating a new group with name:", name)

	if name == "" {
		return errors.New("subject name is required")
	}

	if teacherId == "" {
		return errors.New("teacher is required")
	}

	subjectId := uuid.New().String()
	newSubject := subjectsModels.Subject{
		Id:        subjectId,
		Name:      name,
		TeacherId: teacherId,
	}

	if err := s.subjectsRepository.Create(c, newSubject); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Subject with name '%s' already exists", newSubject.Name)
				return errorsModels.ErrSubjectExists
			default:
				log.Printf("Error creating subject '%s' in the repository: %s", newSubject.Name, err.Error())
				return errorsModels.ErrServer
			}
		}
	}

	subjectTeacherId := uuid.New().String()

	if err := s.subjectsRepository.AddTeacherToSubject(c, subjectTeacherId, subjectId, teacherId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Subject with id %s already has teacher with id %s", subjectId, teacherId)
				return errorsModels.ErrSubjectWithThisTeacherExists
			default:
				log.Printf("Error add teacher with id %s to subject with id %s", teacherId, subjectId)
				return errorsModels.ErrServer
			}
		}
	}

	return nil
}

func (s *service) GetAll(c *fiber.Ctx) ([]subjectsModels.Subject, error) {
	log.Println("Retrieving all subjects")
	subjectsRepo, err := s.subjectsRepository.GetAll(c)

	subjects := subjectsMappers.FromSubjectsRepoModelToEntity(subjectsRepo)

	return subjects, err
}

func (s *service) AddSubjectToGroup(c *fiber.Ctx, subjectGroup subjectsModels.SubjectGroup) error {
	log.Printf("Attempting to add subject with ID '%s' and teacher with ID '%s' to group '%s'.",
		subjectGroup.SubjectId, subjectGroup.TeacherId, subjectGroup.GroupId)

	if subjectGroup.GroupId == "" {
		log.Printf("Validation error: group_id is empty.")
		return errors.New("group_id is required")
	}

	if subjectGroup.TeacherId == "" {
		log.Printf("Validation error: teacher_id is empty.")
		return errors.New("teacher_id is required")
	}

	subject, err := s.subjectsRepository.GetById(c, subjectGroup.SubjectId)
	if err != nil {
		log.Printf("Error retrieving subject '%s': %v", subjectGroup.SubjectId, err)
		return err
	}

	if subject == nil {
		log.Printf("Subject '%s' does not exist.", subjectGroup.SubjectId)
		return errorsModels.ErrSubjectDoesNotExist
	}

	groupExist, err := s.groupsRepository.CheckGroupExists(c, subjectGroup.GroupId)
	if err != nil {
		log.Printf("Error checking existence of group '%s': %v", subjectGroup.GroupId, err)
		return err
	}

	if !groupExist {
		log.Printf("Group '%s' does not exist.", subjectGroup.GroupId)
		return errorsModels.ErrGroupDoesNotExist
	}

	teacherExist, err := s.usersRepository.GetById(c, subjectGroup.TeacherId)

	if err != nil {
		log.Printf("Error checking existence of teacher '%s': %v", subjectGroup.TeacherId, err)
		return err
	}

	if teacherExist == nil {
		log.Printf("Teacher '%s' does not exist", subjectGroup.TeacherId)
		return errorsModels.ErrUserDoesNotExist
	}

	if teacherExist.Role != usersModels.Teacher {
		log.Printf("User with ID '%s' is not teacher", subjectGroup.TeacherId)
		return errorsModels.ErrNoPermission
	}

	if err := s.subjectsRepository.AddSubjectToGroup(c, subjectGroup); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Subject '%s' already exists in group '%s'.", subjectGroup.SubjectId, subjectGroup.GroupId)
				return errorsModels.ErrSubjectExists
			default:
				log.Printf("Postgres error while adding subject '%s' to group '%s': %v",
					subjectGroup.SubjectId, subjectGroup.GroupId, err)
				return errorsModels.ErrServer
			}
		}
	}

	log.Printf("Subject '%s' successfully added to group '%s'.", subjectGroup.SubjectId, subjectGroup.GroupId)
	return nil
}

func (s *service) DeleteSubjectFromGroup(c *fiber.Ctx, subjectId string, groupId string) error {
	log.Printf("Attempting to delete subject '%s' from group '%s'.", subjectId, groupId)

	if subjectId == "" {
		log.Println("Validation error: 'subject_id' is empty.")
		return errors.New("subject_id is required")
	}

	subject, err := s.subjectsRepository.GetById(c, subjectId)
	if err != nil {
		log.Printf("Error retrieving subject '%s': %v", subjectId, err)
		return err
	}
	if subject == nil {
		log.Printf("Subject '%s' does not exist.", subjectId)
		return errorsModels.ErrSubjectDoesNotExist
	}

	subjectInGroup, err := s.subjectsRepository.IsSubjectInGroup(c, subjectId, groupId)
	if err != nil {
		log.Printf("Error checking if group '%s' has subject '%s': %v", groupId, subjectId, err)
		return err
	}

	if !subjectInGroup {
		log.Printf("Group '%s' has not subject '%s'.", groupId, subjectId)
		return errorsModels.ErrGroupNotHasSubject
	}

	if err := s.subjectsRepository.DeleteSubjectFromGroup(c, subjectId, groupId); err != nil {
		log.Printf("Error deleting subject '%s' from group '%s': %v", subjectId, groupId, err)
		return err
	}

	log.Printf("Subject '%s' successfully deleted from group '%s'.", subjectId, groupId)
	return nil
}

func (s *service) DeleteSubject(c *fiber.Ctx, subjectId string) error {
	log.Printf("Attempting to delete subject '%s'.", subjectId)

	subject, err := s.subjectsRepository.GetById(c, subjectId)
	if err != nil {
		log.Printf("Error retrieving subject '%s': %v", subjectId, err)
		return err
	}

	if subject == nil {
		log.Printf("Subject '%s' does not exist.", subjectId)
		return errorsModels.ErrSubjectDoesNotExist
	}

	if err := s.subjectsRepository.DeleteSubject(c, subjectId); err != nil {
		log.Printf("Error deleting subject '%s': %v", subjectId, err)
		return err
	}

	log.Printf("Subject '%s' successfully deleted.", subjectId)
	return nil
}

func (s *service) UpdateSubject(c *fiber.Ctx, subjectId string, subject subjectsModels.Subject) error {
	log.Printf("Attempting to update subject '%s'.", subjectId)

	existedSubject, err := s.subjectsRepository.GetById(c, subjectId)
	if err != nil {
		log.Printf("Error retrieving subject '%s': %v", subjectId, err)
		return err
	}

	if existedSubject == nil || !existedSubject.Id.Valid {
		log.Printf("Subject '%s' does not exist.", subjectId)
		return errorsModels.ErrSubjectDoesNotExist
	}

	subject.Id = existedSubject.Id.String
	if subject.Name == "" {
		subject.Name = existedSubject.Name.String
	}

	if err := s.subjectsRepository.UpdateSubject(c, subject); err != nil {
		log.Printf("Error updating subject '%s': %v", subjectId, err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Subject with id '%s' already has teacher with id %s", subjectId, subject.TeacherId)
				return errorsModels.ErrSubjectWithThisTeacherExists
			default:
				log.Printf("Error updating subject '%s' in the repository: %s", subjectId, err.Error())
				return errorsModels.ErrServer
			}
		}
		return err
	}

	//subjectTeacherId := uuid.New().String()
	//
	//if err := s.subjectsRepository.AddTeacherToSubject(c, subjectTeacherId, subjectId, subject.TeacherId); err != nil {
	//	var pgErr *pgconn.PgError
	//	if errors.As(err, &pgErr) {
	//		switch pgErr.Code {
	//		case errorsModels.NeedUniqueValueErrCode:
	//			log.Printf("Subject with id %s already has teacher with id %s", subjectId, subject.TeacherId)
	//			return errorsModels.ErrSubjectWithThisTeacherExists
	//		default:
	//			log.Printf("Error add teacher with id %s to subject with id %s", subject.TeacherId, subjectId)
	//			return errorsModels.ErrServer
	//		}
	//	}
	//}

	log.Printf("Subject '%s' successfully updated.", subjectId)
	return nil
}

func (s *service) GetAvailableTeachers(c *fiber.Ctx, teacherName string) ([]usersModels.UserInfo, error) {
	log.Printf("Attempting to get available teachers with name %s", teacherName)
	usersRepo, err := s.usersRepository.GetAvailableTeachers(c, teacherName)
	if err != nil {
		log.Printf("Error get available teachers with name %s", teacherName)
		return nil, err
	}

	users := usersMapper.UsersInfoFromRepoToService(usersRepo)

	return users, nil
}

func (s *service) GetByIdWithDetails(c *fiber.Ctx, subjectId string) (*subjectsModels.SubjectInfo, error) {
	log.Printf("Attempting to get subject with details")
	subjectRepo, err := s.subjectsRepository.GetByIdWithDetails(c, subjectId)
	if err != nil {
		log.Printf("Error retrieving subject with ID: %s, error: %v", subjectId, err)
		return nil, err
	}

	if subjectRepo == nil {
		log.Printf("Subject with ID: %s does not exist", subjectId)
		return nil, errorsModels.ErrSubjectDoesNotExist
	}

	subject := subjectsMappers.FromSubjectDetailsRepoModelToEntity(*subjectRepo)

	log.Printf("Successfully retrieved subject with details")

	return &subject, nil
}
