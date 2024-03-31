package groups

import (
	"errors"
	"log"
	"sad/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	groupsMapper "sad/internal/mappers/groups"
	errorsModels "sad/internal/models/errors"
	groupsModels "sad/internal/models/groups"
)

type service struct {
	groupsRepository repositories.GroupsRepository
	usersRepository  repositories.UserRepository
}

func NewService(groupsRepository repositories.GroupsRepository, usersRepository repositories.UserRepository) *service {
	return &service{
		groupsRepository: groupsRepository,
		usersRepository:  usersRepository,
	}
}

func (s *service) Create(c *fiber.Ctx, number string) error {
	log.Println("Creating a new group with number:", number)

	if number == "" {
		return errors.New("group number is required")
	}

	newGroup := groupsModels.Group{
		Id:     uuid.New().String(),
		Number: number,
	}

	if err := s.groupsRepository.Create(c, newGroup); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("Group with number '%s' already exists", newGroup.Number)
				return errorsModels.ErrGroupExists
			default:
				log.Printf("Error creating group '%s' in the repository: %s", newGroup.Number, err.Error())
				return errorsModels.ErrServer
			}
		} else {
			log.Printf("Error creating group '%s' in the repository: %s", newGroup.Number, err.Error())
			return errorsModels.ErrServer
		}
	}

	return nil
}

func (s *service) GetById(c *fiber.Ctx, groupId string) (*groupsModels.Group, error) {
	log.Printf("Retrieving group with ID: %s\n", groupId)

	groupRepo, err := s.groupsRepository.GetById(c, groupId)
	if err != nil {
		log.Printf("Error retrieving group with ID: %s, error: %v\n", groupId, err)
		return nil, err
	}

	if groupRepo == nil {
		log.Printf("Group with ID: %s does not exist\n", groupId)
		return nil, errorsModels.ErrGroupDoesNotExist
	}

	log.Printf("Successfully retrieved group with users for group ID: %s\n", groupId)

	group := groupsMapper.FromGroupRepoModelToEntity(*groupRepo)

	return &group, nil
}

func (s *service) GetByIdWithUsers(c *fiber.Ctx, groupId string) (*groupsModels.GroupWithUsers, error) {
	log.Printf("Retrieving group with ID: %s\n", groupId)

	groupRepo, err := s.groupsRepository.GetByIdWithUsers(c, groupId)
	if err != nil {
		log.Printf("Error retrieving group with ID: %s, error: %v\n", groupId, err)
		return nil, err
	}

	if groupRepo == nil {
		log.Printf("Group with ID: %s does not exist\n", groupId)
		return nil, errorsModels.ErrGroupDoesNotExist
	}

	group := groupsMapper.FromGroupWithUsersRepoModelToEntity(*groupRepo)

	log.Printf("Successfully retrieved group with users for group ID: %s\n", groupId)

	return &group, nil
}

func (s *service) GetAll(c *fiber.Ctx) ([]groupsModels.Group, error) {
	log.Println("Retrieving all groups")
	groupsRepo, err := s.groupsRepository.GetAll(c)

	groups := groupsMapper.FromGroupsRepoModelToEntity(groupsRepo)

	return groups, err
}

func (s *service) AddUserToGroup(c *fiber.Ctx, groupId string, userId string) error {
	log.Printf("Attempting to add user with ID '%s' to group '%s'.", userId, groupId)

	if userId == "" {
		log.Printf("Validation error: user_id is empty.")
		return errors.New("user_id is required")
	}

	group, err := s.groupsRepository.GetById(c, groupId)
	if err != nil {
		log.Printf("Error retrieving group '%s': %v", groupId, err)
		return err
	}

	if group == nil {
		log.Printf("Group '%s' does not exist.", groupId)
		return errorsModels.ErrGroupDoesNotExist
	}

	userExist, err := s.usersRepository.CheckUserExists(c, userId)
	if err != nil {
		log.Printf("Error checking existence of user '%s': %v", userId, err)
		return err
	}

	if !userExist {
		log.Printf("User '%s' does not exist.", userId)
		return errorsModels.ErrUserDoesNotExist
	}

	if err := s.groupsRepository.AddUserToGroup(c, groupId, userId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case errorsModels.NeedUniqueValueErrCode:
				log.Printf("User '%s' already exists in group '%s'.", userId, groupId)
				return errorsModels.ErrUserExists
			default:
				log.Printf("Postgres error while adding user '%s' to group '%s': %v", userId, groupId, err)
				return errorsModels.ErrServer
			}
		} else {
			log.Printf("Unknown error while adding user '%s' to group '%s': %v", userId, groupId, err)
			return errorsModels.ErrServer
		}
	}

	log.Printf("User '%s' successfully added to group '%s'.", userId, groupId)
	return nil
}

func (s *service) DeleteUserFromGroup(c *fiber.Ctx, groupId string, userId string) error {
	log.Printf("Attempting to delete user '%s' from group '%s'.", userId, groupId)

	if userId == "" {
		log.Println("Validation error: 'userId' is empty.")
		return errors.New("user_id is required")
	}

	group, err := s.groupsRepository.GetById(c, groupId)
	if err != nil {
		log.Printf("Error retrieving group '%s': %v", groupId, err)
		return err
	}
	if group == nil {
		log.Printf("Group '%s' does not exist.", groupId)
		return errorsModels.ErrGroupDoesNotExist
	}

	userInGroup, err := s.groupsRepository.IsUserInGroup(c, groupId, userId)
	if err != nil {
		log.Printf("Error checking if user '%s' is in group '%s': %v", userId, groupId, err)
		return err
	}

	if !userInGroup {
		log.Printf("User '%s' is not in group '%s'.", userId, groupId)
		return errorsModels.ErrUserNotInGroup
	}

	if err := s.groupsRepository.DeleteUserFromGroup(c, groupId, userId); err != nil {
		log.Printf("Error deleting user '%s' from group '%s': %v", userId, groupId, err)
		return err
	}

	log.Printf("User '%s' successfully deleted from group '%s'.", userId, groupId)
	return nil
}

func (s *service) DeleteGroup(c *fiber.Ctx, groupId string) error {
	log.Printf("Attempting to delete group '%s'.", groupId)

	group, err := s.groupsRepository.GetById(c, groupId)
	if err != nil {
		log.Printf("Error retrieving group '%s': %v", groupId, err)
		return err
	}

	if group == nil {
		log.Printf("Group '%s' does not exist.", groupId)
		return errorsModels.ErrGroupDoesNotExist
	}

	if err := s.groupsRepository.DeleteGroup(c, groupId); err != nil {
		log.Printf("Error deleting group '%s': %v", groupId, err)
		return err
	}

	log.Printf("Group '%s' successfully deleted.", groupId)
	return nil
}

func (s *service) UpdateGroup(c *fiber.Ctx, groupId string, group groupsModels.Group) error {
	log.Printf("Attempting to update group '%s'.", groupId)

	existedGroup, err := s.groupsRepository.GetById(c, groupId)
	if err != nil {
		log.Printf("Error retrieving group '%s': %v", groupId, err)
		return err
	}

	if existedGroup == nil || !existedGroup.Id.Valid {
		log.Printf("Group '%s' does not exist.", groupId)
		return errorsModels.ErrGroupDoesNotExist
	}

	group.Id = existedGroup.Id.String
	if group.Number == "" {
		group.Number = existedGroup.Number.String
	}

	if err := s.groupsRepository.UpdateGroup(c, group); err != nil {
		log.Printf("Error updating group '%s': %v", groupId, err)
		return err
	}

	log.Printf("Group '%s' successfully updated.", groupId)
	return nil
}
