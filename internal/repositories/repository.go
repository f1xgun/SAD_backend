package repositories

import (
	userModels "sad/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

type AuthRepository interface {
	GetByUUID(c *fiber.Ctx, userUUID string) (*userModels.UserCredentials, error)
	GetByLogin(c *fiber.Ctx, login string) (*userModels.UserCredentials, error)
	Create(c *fiber.Ctx, user userModels.User) error
}
