package users

import (
	"errors"
	"log"
	usersMapper "sad/internal/mappers/users"
	errorsModels "sad/internal/models/errors"
	usersModels "sad/internal/models/users"
	"sad/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type service struct {
	userRepository repositories.UserRepository
}

func NewService(userRepository repositories.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) EditUser(c *fiber.Ctx, userId string, newRole usersModels.UserRole, newName string) error {
	if len(newName) == 0 {
		log.Printf("Error chagne user info: user's name should be no empty")
		return errors.New("user's name should be no empty")
	}

	adminID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("Failed to assert type for userID from Locals")
		return errorsModels.ErrServer
	}
	log.Printf("Admin with ID %s is attempting to change the info of user %s", adminID, userId)

	if adminID == userId {
		log.Printf("Admin with ID %s attempted to change their own info", adminID)
		return errorsModels.ErrChangeOwnRole
	}

	exists, err := s.userRepository.CheckUserExists(c, userId)
	if err != nil {
		log.Printf("Error checking if user exists for userId: %s, error: %v", userId, err)
		return err
	}
	if !exists {
		log.Printf("User with userId: %s does not exist", userId)
		return errorsModels.ErrUserNotFound
	}

	err = s.userRepository.ChangeUserInfo(c, userId, newRole, newName)

	if err != nil {
		log.Printf("Error changing user info: %v", err)
	} else {
		log.Printf("User info changed successfully for user %s by admin %s", userId, adminID)
	}

	return err
}

func (s *service) CheckIsUserRoleAllowed(c *fiber.Ctx, allowedRoles []usersModels.UserRole, userId string) (bool, error) {
	user, err := s.userRepository.GetById(c, userId)
	if err != nil {
		log.Printf("Error fetching user data: %v", err)
		return false, err
	}

	for _, allowedRole := range allowedRoles {
		if user.Role == allowedRole {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) GetUserInfo(c *fiber.Ctx, userId string) (*usersModels.UserInfo, error) {
	userRepoInfo, err := s.userRepository.GetUserInfo(c, userId)

	if err != nil {
		log.Printf("Error fetching user info: %v", err)
		return nil, err
	}

	user := usersMapper.UserInfoFromRepoToService(*userRepoInfo)

	return &user, nil
}

func (s *service) GetUsersInfo(c *fiber.Ctx) ([]usersModels.UserInfo, error) {
	usersRepoInfo, err := s.userRepository.GetUsersInfo(c)

	if err != nil {
		log.Printf("Error fetching users info: %v", err)
		return nil, err
	}

	usersInfo := usersMapper.UsersInfoFromRepoToService(usersRepoInfo)

	return usersInfo, nil
}

func (s *service) DeleteUser(c *fiber.Ctx, userId string) error {
	err := s.userRepository.DeleteUser(c, userId)

	if err != nil {
		log.Printf("Error deleting user with id %s, err %v", userId, err)
		return err
	}

	return nil
}
