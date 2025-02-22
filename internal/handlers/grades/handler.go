package grades

import (
	"errors"
	"log/slog"
	"net/http"
	errorsModels "sad/internal/models/errors"
	"sad/internal/models/grades"
	"sad/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type GradesHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAllStudentGrades(w http.ResponseWriter, r *http.Request)
	GetStudentGradesBySubjectAndGroup(w http.ResponseWriter, r *http.Request)
	GetGradesInCsvReport(w http.ResponseWriter, r *http.Request)
}

type gradesHandler struct {
	gradesService services.GradesService
	log           *slog.Logger
}

func NewGradesHandler(logger *slog.Logger, gradesService services.GradesService) GradesHandler {
	return &gradesHandler{
		gradesService: gradesService,
		log:           logger,
	}
}

func (h *gradesHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.grades.Create"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var body grades.Grade

	if err := render.DecodeJSON(r.Body, &body); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("Attempting to create grade", slog.Any("grade", body))

	err := h.gradesService.Create(r.Context(), body)
	if err != nil {
		var status int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("Failed to create grade: server error", slog.String("error", err.Error()))
		default:
			status = http.StatusBadRequest
			errMsg = err.Error()
			log.Error("Failed to create grade: bad request", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Grade created successfully", slog.Any("grade", body))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Grade created successfully"})
}

func (h *gradesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.grades.Delete"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	gradeId := chi.URLParam(r, "grade_id")
	if gradeId == "" {
		log.Error("grade_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "grade_id is required"})
		return
	}

	log.Info("Attempting to delete grade", slog.String("grade_id", gradeId))

	err := h.gradesService.Delete(r.Context(), gradeId)
	if err != nil {
		var status int

		switch {
		case errors.Is(err, errorsModels.ErrGradeDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Grade not found", slog.String("grade_id", gradeId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to delete grade", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Grade deleted successfully", slog.String("grade_id", gradeId))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Grade deleted successfully"})
}

func (h *gradesHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.grades.Update"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	gradeId := chi.URLParam(r, "grade_id")
	if gradeId == "" {
		log.Error("grade_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "grade_id is required"})
		return
	}

	var grade grades.Grade
	if err := render.DecodeJSON(r.Body, &grade); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("Attempting to update grade", slog.String("grade_id", gradeId))

	evaluation := grade.Evaluation
	comment := grade.Comment
	err := h.gradesService.Update(r.Context(), gradeId, &evaluation, comment)
	if err != nil {
		var status int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrGradeExists):
			status = http.StatusConflict
			errMsg = "Grade with this user_id and subject_id already exists"
			log.Warn("Grade already exists", slog.String("grade_id", gradeId))
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("Server error", slog.String("error", err.Error()))
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("Failed to update grade", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Grade updated successfully", slog.String("grade_id", gradeId))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Grade updated successfully"})
}

func (h *gradesHandler) GetAllStudentGrades(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.grades.GetAllStudentGrades"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	studentId := chi.URLParam(r, "student_id")
	if studentId == "" {
		log.Error("student_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "student_id is required"})
		return
	}

	var subjectId *string
	if subject := r.URL.Query().Get("subject_id"); subject != "" {
		subjectId = &subject
	}

	isFinal := false
	if r.URL.Query().Get("is_final") == "true" {
		isFinal = true
	}

	log.Info("Fetching all grades for student",
		slog.String("student_id", studentId),
		slog.Bool("is_final", isFinal),
		slog.Any("subject_id", subjectId),
	)

	grades, err := h.gradesService.GetAllStudentGrades(r.Context(), studentId, isFinal, subjectId)
	if err != nil {
		var status int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrUserDoesNotExist):
			status = http.StatusNotFound
			errMsg = "User with this student_id doesn't exist"
			log.Warn("User not found", slog.String("student_id", studentId))
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("Server error", slog.String("error", err.Error()))
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("Failed to get user grades", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Successfully fetched grades for student", slog.String("student_id", studentId))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, grades)
}

func (h *gradesHandler) GetStudentGradesBySubjectAndGroup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.grades.GetStudentGradesBySubjectAndGroup"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := r.URL.Query().Get("group_id")
	subjectId := r.URL.Query().Get("subject_id")
	isFinalQuery := r.URL.Query().Get("is_final")

	var isFinal *bool
	if isFinalQuery == "true" {
		value := true
		isFinal = &value
	} else if isFinalQuery == "false" {
		value := false
		isFinal = &value
	}

	log.Info("Fetching student grades by subject and group",
		slog.String("group_id", groupId),
		slog.String("subject_id", subjectId),
		slog.Any("is_final", isFinal),
	)

	usersWithGrades, err := h.gradesService.GetStudentsGradesBySubjectAndGroup(r.Context(), subjectId, groupId, isFinal)
	if err != nil {
		log.Error("Failed to fetch student grades", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Successfully fetched student grades by subject and group",
		slog.String("group_id", groupId),
		slog.String("subject_id", subjectId),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, usersWithGrades)
}

func (h *gradesHandler) GetGradesInCsvReport(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.grades.GetGradesInCsvReport"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Generating CSV report for grades")

	report, err := h.gradesService.GetGradesInCsv(r.Context())
	if err != nil {
		log.Error("Failed to generate CSV report", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"grades_report.csv\"")

	log.Info("CSV report generated successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(report))
}
