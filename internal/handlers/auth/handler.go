package auth

import (
	"net/http"
	authModels "sad/internal/models/auth"
	errorsModels "sad/internal/models/errors"
	userModels "sad/internal/models/user"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var user userModels.UserCredentials

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.authService.Login(c, user)

	if err != nil {
		var statusCode int
		var errMsg string

		switch err {
		case errorsModels.ErrUserNotFound:
			statusCode = http.StatusNotFound
			errMsg = "User with this login does not exist"
		case errorsModels.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
			errMsg = "Invalid credentials"
		case errorsModels.ErrServer:
			statusCode = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = "An unexpected error occurred"
		}

		return c.Status(statusCode).JSON(fiber.Map{"error": errMsg})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var user authModels.UserRegistrationRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	uuid, err := h.authService.Register(c, user)

	if err != nil {
		var statusCode int
		var errMsg string

		switch err {
		case errorsModels.ErrUserExists:
			statusCode = http.StatusConflict
			errMsg = "User with this login already exist"
		case errorsModels.ErrServer:
			statusCode = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			statusCode = http.StatusInternalServerError
			errMsg = err.Error()
		}
		return c.Status(statusCode).JSON(fiber.Map{"error": errMsg})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"uuid": uuid})
}
