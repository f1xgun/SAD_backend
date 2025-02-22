package subjects

import (
	"errors"
	"log/slog"
	"net/http"
	errorsModels "sad/internal/models/errors"
	subjectsModels "sad/internal/models/subjects"
	"sad/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	AddSubjectToGroup(w http.ResponseWriter, r *http.Request)
	DeleteSubjectFromGroup(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAvailableTeachers(w http.ResponseWriter, r *http.Request)
	GetWithDetails(w http.ResponseWriter, r *http.Request)
	GetSubjectsByTeacherId(w http.ResponseWriter, r *http.Request)
	GetNewAvailableSubjectsForTeacher(w http.ResponseWriter, r *http.Request)
	EditTeacherSubjects(w http.ResponseWriter, r *http.Request)
}

type subjectsHandler struct {
	subjectsService services.SubjectsService
	log             *slog.Logger
}

func NewSubjectsHandler(logger *slog.Logger, subjectsService services.SubjectsService) Handler {
	return &subjectsHandler{
		subjectsService: subjectsService,
		log:             logger,
	}
}

func (h *subjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.Create"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var body subjectsModels.Subject
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		log.Warn("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("Attempting to create subject", slog.String("subject_name", body.Name))

	err := h.subjectsService.Create(r.Context(), body.Name)
	if err != nil {
		var status int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrSubjectExists):
			status = http.StatusConflict
			errMsg = "Subject with this name already exist"
			log.Warn("Subject already exist", slog.String("subject_name", body.Name))
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("Server error", slog.String("error", err.Error()))
		default:
			status = http.StatusBadRequest
			errMsg = err.Error()
			log.Error("Failed to create subject", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Subject created successfully", slog.String("subject_name", body.Name))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Subject created successfully"})
}

func (h *subjectsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "handler.groups.GetAll"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Fetching all subjects")

	subjects, err := h.subjectsService.GetAll(r.Context())
	if err != nil {
		log.Error("Failed to retrieve subjects: %v", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Successfully fetched all subjects")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, subjects)
}

func (h *subjectsHandler) AddSubjectToGroup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.AddSubjectToGroup"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var body subjectsModels.SubjectGroup
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	body.SubjectId = chi.URLParam(r, "subject_id")

	log.Info(
		"Attemting to add subject to group",
		slog.String("subject_id", body.SubjectId),
		slog.String("group_id", body.GroupId),
	)

	err := h.subjectsService.AddSubjectToGroup(r.Context(), body)

	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
			log.Warn(
				"Subject does not exist",
				slog.String("subject_id", body.SubjectId),
			)
		case errors.Is(err, errorsModels.ErrSubjectExists):
			status = http.StatusConflict
			log.Warn(
				"Subject already exists in group",
				slog.String("subject_id", body.SubjectId),
			)
		//case errors.Is(err, errorsModels.ErrSubjectWithThisTeacherExists):
		//	status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to add subject to group", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info(
		"Subject added to group successfully",
		slog.String("subject_id", body.SubjectId),
		slog.String("group_id", body.GroupId),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Subject added to group successfully"})
}

func (h *subjectsHandler) DeleteSubjectFromGroup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.DeleteSubjectFromGroup"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var body subjectsModels.SubjectGroup
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		log.Error("Failed to decode request body", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	subjectId := chi.URLParam(r, "subject_id")
	groupId := body.GroupId

	log.Info(
		"Attempting to delete subject from group",
		slog.String("subject_id", subjectId),
		slog.String("group_id", groupId),
	)

	err := h.subjectsService.DeleteSubjectFromGroup(r.Context(), subjectId, groupId)

	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Subject not found", slog.String("subject_id", subjectId))
		case errors.Is(err, errorsModels.ErrGroupNotHasSubject):
			status = http.StatusConflict
			log.Warn(
				"Group does not have subject",
				slog.String("subject_id", subjectId),
				slog.String("group_id", groupId),
			)
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to delete subject from group", slog.String("err", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info(
		"Subject deleted from group successfully",
		slog.String("subject_id", subjectId),
		slog.String("group_id", groupId),
	)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Subject deleted from group successfully"})

}

func (h *subjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.Delete"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	subjectId := chi.URLParam(r, "subject_id")

	if subjectId == "" {
		log.Error("subject_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "subject_id is required"})
		return
	}

	log.Info("Attempting to delete subject", slog.String("subject_id", subjectId))

	err := h.subjectsService.DeleteSubject(r.Context(), subjectId)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Subject not found", slog.String("subject_id", subjectId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to delete subject", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	return
}

func (h *subjectsHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.Update"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	subjectId := chi.URLParam(r, "subject_id")
	if subjectId == "" {
		log.Error("subject_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "subject_id is required"})
		return
	}

	var subject subjectsModels.Subject
	if err := render.DecodeJSON(r.Body, &subject); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("Attemting to update subject", slog.String("subject_id", subjectId))

	err := h.subjectsService.UpdateSubject(r.Context(), subjectId, subject)
	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrSubjectExists):
			status = http.StatusConflict
			errMsg = "Subject with this name already exist"
			log.Warn(
				errMsg,
				slog.String("subject_name", subject.Name),
			)
		case errors.Is(err, errorsModels.ErrSubjectWithThisTeacherExists):
			status = http.StatusConflict
			errMsg = "Subject with this teacher already exist"
			log.Warn(errMsg, slog.String("subject_id", subjectId))
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error(errMsg, slog.String("error", err.Error()))
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("Failed to update subject", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Subject updated successfully", slog.String("subject_id", subjectId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Subject update successfully"})
}

func (h *subjectsHandler) GetAvailableTeachers(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.GetAvailableTeachers"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	teacherName := r.URL.Query().Get("name")

	if teacherName == "" {
		log.Error("name is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "name is requried"})
		return
	}

	log.Info("Attemting to get available teachers", slog.String("name", teacherName))

	teachers, err := h.subjectsService.GetAvailableTeachers(r.Context(), teacherName)
	if err != nil {
		log.Warn("Failed to retrieve avabilable teachers", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, teachers)
}

func (h *subjectsHandler) GetWithDetails(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.GetWithDetails"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	subjectId := chi.URLParam(r, "subject_id")

	if subjectId == "" {
		log.Warn("subject_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "subject_id is required"})
		return
	}

	log.Info("Attemting to get details of subject", slog.String("subject_id", subjectId))

	subject, err := h.subjectsService.GetByIdWithDetails(r.Context(), subjectId)
	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrSubjectDoesNotExist):
			status = http.StatusNotFound
			errMsg = "Subject does not exist"
			log.Warn(errMsg, slog.String("subject_id", subjectId))
		default:
			status = http.StatusInternalServerError
			errMsg = "Internal server error"
			log.Error("Failed to get details of subject", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Subject details retrieve successfully", slog.String("subject_id", subjectId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, subject)
	return
}

func (h *subjectsHandler) GetSubjectsByTeacherId(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.GetSubjectsByTeacherId"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	teacherId := r.URL.Query().Get("teacher_id")
	if teacherId == "" {
		selfId, ok := r.Context().Value("user_id").(string)
		if !ok || selfId == "" {
			log.Error("Failed to get user_id from context")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid request body"})
			return
		}
		teacherId = selfId
	}

	log.Info("Attemting to get subjects for teacher", slog.String("teacher_id", teacherId))

	subjects, err := h.subjectsService.GetSubjectsByTeacherId(r.Context(), teacherId)
	if err != nil {
		log.Warn("Failed to retrieve subjects by teacher id", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Successfully get subjects by teacher id", slog.String("teacher_id", teacherId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, subjects)
}

func (h *subjectsHandler) GetNewAvailableSubjectsForTeacher(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.GetNewAvailableSubjectsForTeacher"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	teacherID := r.URL.Query().Get("teacher_id")

	if teacherID == "" {
		log.Warn("teacher_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "teacher_id is required"})
		return
	}

	subjects, err := h.subjectsService.GetNewAvailableSubjectsForTeacher(r.Context(), teacherID)
	if err != nil {
		log.Warn("Failed to retrieve subjects by teacher id", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error"})
		return
	}

	log.Info("Successfully get new available subjects for teacher", slog.String("teacher_id", teacherID))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, subjects)
}

func (h *subjectsHandler) EditTeacherSubjects(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.subjects.EditTeacherSubjects"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	teacherID := r.URL.Query().Get("teacher_id")
	if teacherID == "" {
		log.Warn("teacher_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "teacher_id is required"})
		return
	}

	var subjects []subjectsModels.Subject
	if err := render.DecodeJSON(r.Body, &subjects); err != nil {
		log.Warn("Invalid request body")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}

	log.Info("Attempting to edit teacher subjects", slog.String("teacher_id", teacherID))

	err := h.subjectsService.EditTeacherSubjects(r.Context(), teacherID, subjects)
	if err != nil {
		log.Warn(
			"Failed to edit subjects for teacher",
			slog.String("error", err.Error()),
			slog.String("teacher_id", teacherID),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Teacher's subjects updated succesfully", slog.String("teacher_id", teacherID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Teacher's subjects updated succesfully"})
}
