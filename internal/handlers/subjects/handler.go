package subjects

import (
	"errors"
	"log"
	"net/http"
	errorsModels "sad/internal/models/errors"
	subjectsModels "sad/internal/models/subjects"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SubjectsHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	AddSubjectToGroup(c *fiber.Ctx) error
	DeleteSubjectFromGroup(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type subjectsHandler struct {
	subjectsService services.SubjectsService
}

func NewSubjectsHandler(subjectsService services.SubjectsService) SubjectsHandler {
	return &subjectsHandler{
		subjectsService: subjectsService,
	}
}

func (h *subjectsHandler) Create(c *fiber.Ctx) error {
	var body subjectsModels.Subject
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	name := body.Name
	err := h.subjectsService.Create(c, name)
	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrSubjectExists):
			status = http.StatusConflict
			errMsg = "Subject with this name already exist"
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusBadRequest
			errMsg = err.Error()
		}
		log.Printf("Failed to create subject: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Subject created successfully"})
}

func (h *subjectsHandler) GetAll(c *fiber.Ctx) error {
	subjects, err := h.subjectsService.GetAll(c)
	if err != nil {
		log.Printf("Failed to retrieve subjects: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(subjects)
}

func (h *subjectsHandler) AddSubjectToGroup(c *fiber.Ctx) error {
	var body subjectsModels.SubjectGroup
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}
	body.SubjectId = c.Params("subject_id")

	err := h.subjectsService.AddSubjectToGroup(c, body)

	if err != nil {
		log.Printf("Failed to add subject to group: %v", err)
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
		case errors.Is(err, errorsModels.ErrSubjectExists):
			status = http.StatusConflict
		//case errors.Is(err, errorsModels.ErrSubjectWithThisTeacherExists):
		//	status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Subject added to group successfully"})
}

func (h *subjectsHandler) DeleteSubjectFromGroup(c *fiber.Ctx) error {
	var body subjectsModels.SubjectGroup
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}
	subjectId := c.Params("subject_id")
	groupId := body.GroupId

	err := h.subjectsService.DeleteSubjectFromGroup(c, subjectId, groupId)

	if err != nil {
		log.Printf("Failed to delete subject from group: %v", err)
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
		case errors.Is(err, errorsModels.ErrGroupNotHasSubject):
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Subject removed from group successfully"})
}

func (h *subjectsHandler) Delete(c *fiber.Ctx) error {
	subjectId := c.Params("subject_id")
	err := h.subjectsService.DeleteSubject(c, subjectId)
	if err != nil {
		log.Printf("Failed to delete subject: %v", err)
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(http.StatusOK)
}

func (h *subjectsHandler) Update(c *fiber.Ctx) error {
	subjectId := c.Params("subject_id")
	var subject subjectsModels.Subject
	if err := c.BodyParser(&subject); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	err := h.subjectsService.UpdateSubject(c, subjectId, subject)
	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrSubjectExists):
			status = http.StatusConflict
			errMsg = "Subject with this number already exist"
		//case errors.Is(err, errorsModels.ErrSubjectWithThisTeacherExists):
		//	status = http.StatusConflict
		//	errMsg = err.Error()
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Failed to update subject: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Subject update successfully"})
}
