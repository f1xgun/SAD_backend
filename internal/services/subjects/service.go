package subjects

import (
	"errors"
	"log"
	subjectsMappers "sad/internal/mappers/subjects"
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
}

func NewService(groupsRepository repositories.GroupsRepository, subjectsRepository repositories.SubjectsRepository) *service {
	return &service{
		groupsRepository:   groupsRepository,
		subjectsRepository: subjectsRepository,
	}
}

func (s *service) Create(c *fiber.Ctx, name string) error {
	log.Println("Creating a new group with name:", name)

	if name == "" {
		return errors.New("subject name is required")
	}

	newSubject := subjectsModels.Subject{
		Id:   uuid.New().String(),
		Name: name,
	}

	if err := s.subjectsRepository.Create(c, newSubject); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Subject with name '%s' already exists", newSubject.Name)
				return errorsModels.ErrSubjectExists
			default:
				log.Printf("Error creating subject '%s' in the repository: %s", newSubject.Name, err.Error())
				return errorsModels.ErrServer
			}
		} else {
			log.Printf("Error creating subject '%s' in the repository: %s", newSubject.Name, err.Error())
			return errorsModels.ErrServer
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

func (s *service) AddSubjectToGroup(c *fiber.Ctx, subjectId string, groupId string) error {
	log.Printf("Attempting to add subject with ID '%s' to group '%s'.", subjectId, groupId)

	if subjectId == "" {
		log.Printf("Validation error: subject_id is empty.")
		return errors.New("subject_id is required")
	}

	group, err := s.subjectsRepository.GetById(c, subjectId)
	if err != nil {
		log.Printf("Error retrieving subject '%s': %v", subjectId, err)
		return err
	}

	if group == nil {
		log.Printf("Subject '%s' does not exist.", subjectId)
		return errorsModels.ErrSubjectDoesNotExist
	}

	groupExist, err := s.groupsRepository.CheckGroupExists(c, groupId)
	if err != nil {
		log.Printf("Error checking existence of group '%s': %v", groupId, err)
		return err
	}

	if !groupExist {
		log.Printf("Group '%s' does not exist.", groupId)
		return errorsModels.ErrGroupDoesNotExist
	}

	if err := s.subjectsRepository.AddSubjectToGroup(c, subjectId, groupId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Subject '%s' already exists in group '%s'.", subjectId, groupId)
				return errorsModels.ErrSubjectExists
			default:
				log.Printf("Postgres error while adding subject '%s' to group '%s': %v", subjectId, groupId, err)
				return errorsModels.ErrServer
			}
		} else {
			log.Printf("Unknown error while adding subject '%s' to group '%s': %v", subjectId, groupId, err)
			return errorsModels.ErrServer
		}
	}

	log.Printf("Subject '%s' successfully added to group '%s'.", subjectId, groupId)
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
		return err
	}

	log.Printf("Subject '%s' successfully updated.", subjectId)
	return nil
}
