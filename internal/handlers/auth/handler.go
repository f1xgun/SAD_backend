package auth

import (
	"net/http"
	"sad/internal/models/auth"
	"sad/internal/models/user"
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"token": token})
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var user auth.UserRegistrationRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	uuid, err := h.authService.Register(c, user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"uuid": uuid})
}
