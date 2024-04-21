package users

import (
	"log"
	"net/http"
	usersModels "sad/internal/models/users"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

func AllowedRoleMiddleware(usersService services.UserService, allowedRoles []usersModels.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("Check user has allowed role")
		userId, ok := c.Locals("userID").(string)
		if !ok {
			log.Println("Failed to assert type for userID from Locals")
			return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "failed to get user id"})
		}

		userHasAllowedRole, err := usersService.CheckIsUserRoleAllowed(c, allowedRoles, userId)
		if err != nil {
			log.Printf("Error check is user admin: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
		}

		if !userHasAllowedRole {
			log.Printf("No permission: %v", err)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"message": "no permission"})
		}

		return c.Next()
	}
}
