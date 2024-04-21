package grades

import (
	"errors"
	"log"
	"net/http"
	errorsModels "sad/internal/models/errors"
	gradesModels "sad/internal/models/grades"
	"sad/internal/services"

	"github.com/gofiber/fiber/v2"
)

type GradesHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetAllStudentGrades(c *fiber.Ctx) error
}

type gradesHandler struct {
	gradesService services.GradesService
}

func NewGradesHandler(gradesService services.GradesService) GradesHandler {
	return &gradesHandler{
		gradesService: gradesService,
	}
}

func (h *gradesHandler) Create(c *fiber.Ctx) error {
	var body gradesModels.Grade
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	err := h.gradesService.Create(c, body)
	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusBadRequest
			errMsg = err.Error()
		}
		log.Printf("Failed to create grade: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Grade created successfully"})
}

func (h *gradesHandler) Delete(c *fiber.Ctx) error {
	gradeId := c.Params("grade_id")
	err := h.gradesService.Delete(c, gradeId)
	if err != nil {
		log.Printf("Failed to delete grade: %v", err)
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrGradeDoesNotExist):
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(http.StatusOK)
}

func (h *gradesHandler) Update(c *fiber.Ctx) error {
	gradeId := c.Params("grade_id")
	var grade gradesModels.Grade
	if err := c.BodyParser(&grade); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"error": "invalid request body"})
	}

	evaluation := grade.Evaluation

	if err := h.gradesService.Update(c, gradeId, evaluation); err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrGradeExists):
			status = http.StatusConflict
			errMsg = "Grade with this user_id and subject_id already exist"
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Failed to update grade: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Grade update successfully"})
}

func (h *gradesHandler) GetAllStudentGrades(c *fiber.Ctx) error {
	studentId := c.Params("student_id")

	grades, err := h.gradesService.GetAllStudentGrades(c, studentId)

	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrUserDoesNotExist):
			status = http.StatusNotFound
			errMsg = "User with this student_id doesn't exist"
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
		}
		log.Printf("Failed to get user grades: %v", errMsg)
		return c.Status(status).JSON(fiber.Map{"error": errMsg})
	}
	return c.Status(http.StatusOK).JSON(grades)
}
