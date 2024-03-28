package users

import (
	"log"
	"net/http"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(usersService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIsAdmin, err := usersService.CheckUserIsAdmin(c)
		if err != nil {
			log.Printf("Error check is user admin: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
		}

		if !userIsAdmin {
			log.Printf("No admin permission: %v", err)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"message": "no permission"})
		}

		return c.Next()
	}
}
