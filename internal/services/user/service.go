package user

import (
	"log"
	errorsModels "sad/internal/models/errors"
	userModels "sad/internal/models/user"
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

func (s *service) EditRole(c *fiber.Ctx, userId string, newRole userModels.UserRole) error {
	adminID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("Failed to assert type for userID from Locals")
		return errorsModels.ErrServer
	}
	log.Printf("Admin with ID %s is attempting to change the role of user %s", adminID, userId)

	if adminID == userId {
		log.Printf("Admin with ID %s attempted to change their own role", adminID)
		return errorsModels.ErrChangeOwnRole
	}

	admin, err := s.userRepository.GetByUUID(c, adminID)
	if err != nil {
		log.Printf("Error fetching admin data: %v", err)
		return err
	}

	if admin.Role != userModels.Admin {
		log.Printf("User with ID %s has no permissions to change roles", adminID)
		return errorsModels.ErrNoPermission
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

	err = s.userRepository.ChangeUserRole(c, userId, newRole)

	if err != nil {
		log.Printf("Error changing user role: %v", err)
	} else {
		log.Printf("User role changed successfully for user %s by admin %s", userId, adminID)
	}

	return err
}
