package user

import (
	"log"
	"net/http"
	errorsModels "sad/internal/models/errors"
	userModels "sad/internal/models/user"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	EditRole(c *fiber.Ctx) error
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) EditRole(c *fiber.Ctx) error {
	userID := c.Query("userId")

	var user userModels.UserCredentials
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	err := h.userService.EditRole(c, userID, user.Role)
	if err != nil {
		var statusCode int
		var errMsg string
		switch err {
		case errorsModels.ErrNoPermission:
			statusCode = http.StatusForbidden
			errMsg = "No permission to change user role"
		case errorsModels.ErrUserNotFound:
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
		case errorsModels.ErrChangeOwnRole:
			statusCode = http.StatusForbidden
			errMsg = "Cannot change own role"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Error occurred while changing user role: %s, Status Code: %d", errMsg, statusCode)
		return c.Status(statusCode).JSON(&fiber.Map{"error": errMsg})
	}

	return c.SendStatus(http.StatusOK)
}
