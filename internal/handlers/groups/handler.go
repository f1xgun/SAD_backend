package groups

import (
	"log"
	"net/http"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"

	errorsModels "sad/internal/models/errors"
	groupsModels "sad/internal/models/groups"
)

type GroupsHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	AddUserToGroup(c *fiber.Ctx) error
	DeleteUserFromGroup(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type groupsHandler struct {
	groupsService services.GroupsService
}

func NewGroupsHandler(groupsService services.GroupsService) GroupsHandler {
	return &groupsHandler{
		groupsService: groupsService,
	}
}

func (h *groupsHandler) Create(c *fiber.Ctx) error {
	var body groupsModels.Group
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	number := body.Number
	err := h.groupsService.Create(c, number)
	if err != nil {
		var status int
		var errMsg string
		switch err {
		case errorsModels.ErrGroupExists:
			status = http.StatusConflict
			errMsg = "Group with this number already exist"
		case errorsModels.ErrServer:
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusBadRequest
			errMsg = err.Error()
		}
		log.Printf("Failed to create group: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Group created successfully"})
}

func (h *groupsHandler) Get(c *fiber.Ctx) error {
	groupId := c.Params("group_id")
	group, err := h.groupsService.GetById(c, groupId)
	if err != nil {
		log.Printf("Failed to retrieve group: %v", err)
		var status int
		switch err {
		case errorsModels.ErrGroupDoesNotExist:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(group)
}

func (h *groupsHandler) GetAll(c *fiber.Ctx) error {
	groups, err := h.groupsService.GetAll(c)
	if err != nil {
		log.Printf("Failed to retrieve groups: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
}

func (h *groupsHandler) AddUserToGroup(c *fiber.Ctx) error {
	var body groupsModels.UserGroup
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}
	groupId := c.Params("group_id")
	userId := body.UserId

	err := h.groupsService.AddUserToGroup(c, groupId, userId)

	if err != nil {
		log.Printf("Failed to add user to group: %v", err)
		var status int
		switch err {
		case errorsModels.ErrUserExists, errorsModels.ErrGroupDoesNotExist:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User added to group successfully"})
}

func (h *groupsHandler) DeleteUserFromGroup(c *fiber.Ctx) error {
	var body groupsModels.UserGroup
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}
	groupId := c.Params("group_id")
	userId := body.UserId

	err := h.groupsService.DeleteUserFromGroup(c, groupId, userId)

	if err != nil {
		log.Printf("Failed to delete user from group: %v", err)
		var status int
		switch err {
		case errorsModels.ErrGroupDoesNotExist, errorsModels.ErrUserNotInGroup:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User removed from group successfully"})
}

func (h *groupsHandler) Delete(c *fiber.Ctx) error {
	groupId := c.Params("group_id")
	err := h.groupsService.DeleteGroup(c, groupId)
	if err != nil {
		log.Printf("Failed to delete group: %v", err)
		var status int
		switch err {
		case errorsModels.ErrGroupDoesNotExist:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(http.StatusOK)
}

func (h *groupsHandler) Update(c *fiber.Ctx) error {
	groupId := c.Params("group_id")
	var group groupsModels.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	err := h.groupsService.UpdateGroup(c, groupId, group)
	if err != nil {
		var status int
		var errMsg string
		switch err {
		case errorsModels.ErrGroupExists:
			status = http.StatusConflict
			errMsg = "Group with this number already exist"
		case errorsModels.ErrServer:
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Failed to update group: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Group update successfully"})
}