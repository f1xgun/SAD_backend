package users

import (
	"errors"
	"log"
	"net/http"
	"sad/internal/models/errors"
	"sad/internal/models/users"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	EditUser(c *fiber.Ctx) error
	GetUserInfo(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetUserInfoByToken(c *fiber.Ctx) error
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) EditUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	var user usersModels.UserInfo
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	err := h.userService.EditUser(c, userID, user.Role, user.Name)
	if err != nil {
		var statusCode int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrNoPermission):
			statusCode = http.StatusForbidden
			errMsg = "No permission to change user info"
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
		case errors.Is(err, errorsModels.ErrChangeOwnRole):
			statusCode = http.StatusForbidden
			errMsg = "Cannot change own info"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Error occurred while changing user info: %s, Status Code: %d", errMsg, statusCode)
		return c.Status(statusCode).JSON(&fiber.Map{"error": errMsg})
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) GetUserInfo(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	userInfo, err := h.userService.GetUserInfo(c, userID)
	if err != nil {
		var statusCode int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Error occurred while fetching user info: %s", errMsg)
		return c.Status(statusCode).JSON(&fiber.Map{"error": errMsg})
	}

	return c.Status(http.StatusOK).JSON(userInfo)
}

func (h *userHandler) GetUsers(c *fiber.Ctx) error {
	usersInfo, err := h.userService.GetUsersInfo(c)
	if err != nil {
		log.Printf("Error occurred while fetching users info: %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(usersInfo)
}

func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	err := h.userService.DeleteUser(c, userID)

	if err != nil {
		log.Printf("Error occured while deleting user with id: %s, err: %s", userID, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) GetUserInfoByToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		log.Println("Failed to assert type for userID from Locals")
		return errorsModels.ErrServer
	}

	userInfo, err := h.userService.GetUserInfo(c, userID)
	if err != nil {
		var statusCode int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrUserNotFound):
			statusCode = http.StatusNotFound
			errMsg = "User with this id does not exist"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Error occurred while fetching user info: %s", errMsg)
		return c.Status(statusCode).JSON(&fiber.Map{"error": errMsg})
	}

	return c.Status(http.StatusOK).JSON(userInfo)
}
