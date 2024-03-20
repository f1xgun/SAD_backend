package repository

import (
	user "sad/internal/repository/models/user"

	"github.com/gofiber/fiber/v2"
)

type AuthRepository interface {
	GetByUUID(c *fiber.Ctx, userUUID string) (*user.User, error)
	GetByLogin(c *fiber.Ctx, login string) (*user.User, error)
	Create(c *fiber.Ctx, user *user.User) error
}
