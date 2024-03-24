package repositories

import (
	userModels "sad/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	GetByUUID(c *fiber.Ctx, userId string) (*userModels.UserCredentials, error)
	GetByLogin(c *fiber.Ctx, login string) (*userModels.UserRepoModel, error)
	Create(c *fiber.Ctx, user userModels.User) error
	ChangeUserRole(c *fiber.Ctx, userId string, newRole userModels.UserRole) error
	CheckUserExists(c *fiber.Ctx, userId string) (bool, error)
}
