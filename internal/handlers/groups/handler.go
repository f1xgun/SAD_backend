package groups

import (
	"errors"
	"log/slog"
	"net/http"
	errorsModels "sad/internal/models/errors"
	groupsModels "sad/internal/models/groups"
	"sad/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	GetWithDetails(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	AddUserToGroup(w http.ResponseWriter, r *http.Request)
	DeleteUserFromGroup(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAvailableNewUsers(w http.ResponseWriter, r *http.Request)
	GetGroupsWithSubjectsByTeacher(w http.ResponseWriter, r *http.Request)
	GetTeacherGroupsBySubject(w http.ResponseWriter, r *http.Request)
}

type groupsHandler struct {
	groupsService services.GroupsService
	log           *slog.Logger
}

func NewGroupsHandler(logger *slog.Logger, groupsService services.GroupsService) Handler {
	return &groupsHandler{
		groupsService: groupsService,
		log:           logger,
	}
}

func (h *groupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.Create"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var body groupsModels.Group
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("Attempting to create group", slog.String("group_number", body.Number))

	err := h.groupsService.Create(r.Context(), body.Number)
	if err != nil {
		var status int
		var errMsg string

		switch {
		case errors.Is(err, errorsModels.ErrGroupExists):
			status = http.StatusConflict
			errMsg = "Group with this number already exists"
			log.Warn("Group already exists", slog.String("group_number", body.Number))
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("Server error", slog.String("error", err.Error()))
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("Failed to create group", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Group created successfully", slog.String("group_number", body.Number))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Group created successfully"})
}

func (h *groupsHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.Get"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := chi.URLParam(r, "group_id")
	if groupId == "" {
		log.Error("group_id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "group_id is required"})
		return
	}

	log.Info("Fetching group by ID", slog.String("group_id", groupId))

	group, err := h.groupsService.GetById(r.Context(), groupId)
	if err != nil {
		var status int

		switch {
		case errors.Is(err, errorsModels.ErrGroupDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Group not found", slog.String("group_id", groupId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to retrieve group", slog.String("error", err.Error()))
		}

		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Group retrieved successfully", slog.String("group_id", groupId))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, group)
}

func (h *groupsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.GetAll"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Fetching all groups")

	groups, err := h.groupsService.GetAll(r.Context())
	if err != nil {
		log.Error("Failed to retrieve groups", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Successfully fetched all groups")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, groups)
}

func (h *groupsHandler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.AddUserToGroup"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var body groupsModels.UserGroup
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	groupId := chi.URLParam(r, "group_id")
	userId := body.UserId

	log.Info("Attempting to add user to group", slog.String("group_id", groupId), slog.String("user_id", userId))

	err := h.groupsService.AddUserToGroup(r.Context(), groupId, userId)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrUserExists):
			status = http.StatusConflict
			log.Warn("User already exists in group", slog.String("user_id", userId), slog.String("group_id", groupId))
		case errors.Is(err, errorsModels.ErrGroupDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Group not found", slog.String("group_id", groupId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to add user to group", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("User added to group successfully", slog.String("user_id", userId), slog.String("group_id", groupId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User added to group successfully"})
}

func (h *groupsHandler) DeleteUserFromGroup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.DeleteUserFromGroup"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := chi.URLParam(r, "group_id")
	userId := chi.URLParam(r, "user_id")

	log.Info("Attempting to delete user from group", slog.String("group_id", groupId), slog.String("user_id", userId))

	err := h.groupsService.DeleteUserFromGroup(r.Context(), groupId, userId)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrGroupDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Group not found", slog.String("group_id", groupId))
		case errors.Is(err, errorsModels.ErrUserNotInGroup):
			status = http.StatusNotFound
			log.Warn("User not in group", slog.String("user_id", userId), slog.String("group_id", groupId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to delete user from group", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("User removed from group successfully", slog.String("user_id", userId), slog.String("group_id", groupId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User removed from group successfully"})
}

func (h *groupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.Delete"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := chi.URLParam(r, "group_id")

	log.Info("Attempting to delete group", slog.String("group_id", groupId))

	err := h.groupsService.DeleteGroup(r.Context(), groupId)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrGroupDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Group not found", slog.String("group_id", groupId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to delete group", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Group deleted successfully", slog.String("group_id", groupId))
	w.WriteHeader(http.StatusOK)
}

func (h *groupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.Update"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := chi.URLParam(r, "group_id")
	var group groupsModels.Group

	if err := render.DecodeJSON(r.Body, &group); err != nil {
		log.Error("Failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	log.Info("Attempting to update group", slog.String("group_id", groupId))

	err := h.groupsService.UpdateGroup(r.Context(), groupId, group)
	if err != nil {
		var status int
		var errMsg string
		switch {
		case errors.Is(err, errorsModels.ErrGroupExists):
			status = http.StatusConflict
			errMsg = "Group with this number already exists"
			log.Warn("Group already exists", slog.String("group_id", groupId))
		case errors.Is(err, errorsModels.ErrServer):
			status = http.StatusInternalServerError
			errMsg = "Server error"
			log.Error("Server error", slog.String("error", err.Error()))
		default:
			status = http.StatusInternalServerError
			errMsg = err.Error()
			log.Error("Failed to update group", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": errMsg})
		return
	}

	log.Info("Group updated successfully", slog.String("group_id", groupId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Group updated successfully"})
}

func (h *groupsHandler) GetWithDetails(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.GetWithDetails"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := chi.URLParam(r, "group_id")

	log.Info("Fetching group with details", slog.String("group_id", groupId))

	group, err := h.groupsService.GetWithDetailsById(r.Context(), groupId)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrGroupDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Group not found", slog.String("group_id", groupId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to retrieve group", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Group with details retrieved successfully", slog.String("group_id", groupId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, group)
}

func (h *groupsHandler) GetAvailableNewUsers(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.GetAvailableNewUsers"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	groupId := chi.URLParam(r, "group_id")
	login := r.URL.Query().Get("login")

	log.Info("Fetching available new users for group", slog.String("group_id", groupId), slog.String("login", login))

	users, err := h.groupsService.GetAvailableNewUsers(r.Context(), groupId, login)
	if err != nil {
		var status int
		switch {
		case errors.Is(err, errorsModels.ErrGroupDoesNotExist):
			status = http.StatusNotFound
			log.Warn("Group not found", slog.String("group_id", groupId))
		default:
			status = http.StatusInternalServerError
			log.Error("Failed to retrieve available new users", slog.String("error", err.Error()))
		}
		render.Status(r, status)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Available new users retrieved successfully", slog.String("group_id", groupId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}

func (h *groupsHandler) GetGroupsWithSubjectsByTeacher(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.GetGroupsWithSubjectsByTeacher"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	teacherId := r.URL.Query().Get("teacher_id")

	log.Info("Fetching groups with subjects by teacher", slog.String("teacher_id", teacherId))

	groupsWithSubjects, err := h.groupsService.GetGroupsWithSubjectsByTeacher(r.Context(), teacherId)
	if err != nil {
		log.Error("Failed to retrieve groups with subjects", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Groups with subjects retrieved successfully", slog.String("teacher_id", teacherId))
	render.Status(r, http.StatusOK)
	render.JSON(w, r, groupsWithSubjects)
}

func (h *groupsHandler) GetTeacherGroupsBySubject(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.groups.GetTeacherGroupsBySubject"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.Error("Failed to assert type for userID from context")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": errorsModels.ErrServer.Error()})
		return
	}

	subjectId := r.URL.Query().Get("subject_id")

	log.Info("Fetching teacher groups by subject",
		slog.String("teacher_id", userID),
		slog.String("subject_id", subjectId),
	)

	groups, err := h.groupsService.GetGroupsBySubjectAndTeacher(r.Context(), userID, subjectId)
	if err != nil {
		log.Error("Failed to retrieve teacher groups by subject", slog.String("error", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	log.Info("Teacher groups by subject retrieved successfully",
		slog.String("teacher_id", userID),
		slog.String("subject_id", subjectId),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, groups)
}
