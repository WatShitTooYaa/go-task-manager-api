package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	resp "github.com/WatShitTooYaa/go-task-manager-api/internal/response"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/service"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	// "github.com/WatShitTooYaa/go-task-manager-api/."
)

// func
type DBHandler struct {
	service *service.TaskService
}

func NewDBHandler(service *service.TaskService) *DBHandler {
	return &DBHandler{service: service}
}

func (handler *DBHandler) ListTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	tasks, err := handler.service.GetTasks(ctx)
	if err != nil {
		msg := "Failed to load tasks"

		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)

		resp.InternalError(w, err.Error())
		return
	}
	// respondSuccess(w, tasks, http.StatusOK)
	msg := "Tasks loaded successfully"
	log.Debug().
		Str("request_id", reqID).
		Int("count", len(tasks)).
		Msg(msg)

	resp.SendSuccessResponse(w, "Success get all data", tasks, http.StatusOK)
}

func (handler *DBHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	var input entity.CreateTaskRequest
	// input := CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Err(err).
			Msg("Invalid JSON in create task request")
			// sendResponse(w, err.Error(), false, nil, http.StatusInternalServerError)
		resp.InvalidJSON(w)
		return
	}

	if err := validation.ValidateStruct(input); err != nil {
		log.Warn().
			Str("request_id", reqID).
			Str("validation_error", err.Error()).
			Msg("Validation failed")
		// sendResponse(w, err.Error(), false, nil, http.StatusBadRequest)
		resp.ValidationError(w, err.Error(), nil)
		return
	}

	inputTask := entity.Task{
		Content:   input.Content,
		Completed: false,
		Timestamp: time.Now().Format(time.RFC3339),
		Priority:  input.Priority,
	}

	task, err := handler.service.AddTask(ctx, inputTask)
	if err != nil {
		msg := "Failed to create task"
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)

		resp.InternalError(w, msg)
		return
	}

	log.Info().
		Str("request_id", reqID).
		Uint16("task_id", inputTask.Id).
		Str("content", inputTask.Content).
		Msg("Task created successfully")

	resp.SendSuccessResponse(w, "Task created", task, http.StatusCreated)
}

func (handler *DBHandler) GetSingleTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	idStr := chi.URLParam(r, "id")
	if idStr == "" {

		// respondError(w, "Path must not null", http.StatusBadRequest)
		// sendResponse(w, "Path must not null", false, nil, http.StatusBadRequest)
		resp.BadRequest(w, "Path must not null")
		return
	}

	intId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Str("id", idStr).
			Msg("Invalid ID format")

		// sendResponse(w, "Path must be int", false, nil, http.StatusBadRequest)
		resp.InvalidID(w)
		return
	}

	task, err := handler.service.GetSingleTask(ctx, uint16(intId))
	if err != nil {
		msg := "Task not found"
		// fmt.Println(msg)

		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg(msg)

		resp.TaskNotFound(w, intId)
		// sendResponse(w, msg, false, nil, http.StatusNotFound)

		return
	}

	log.Debug().
		Str("request_id", reqID).
		Int("task_id", intId).
		Msg("Task retrieved successfully")

	// sendResponse(w, "Success", true, task, http.StatusOK)
	resp.SendSuccessResponse(w, "", task, http.StatusOK)
}

func (handler *DBHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	id := chi.URLParam(r, "id")
	if id == "" {
		// sendResponse(w, "Path must not null", false, nil, http.StatusBadRequest)
		resp.BadRequest(w, "Path must not null")

		return
	}

	intId, err := strconv.Atoi(id)

	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Str("id", id).
			Msg("Invalid ID format")

		// sendResponse(w, "Invalid ID format", false, nil, http.StatusBadRequest)
		// BadRequest(w, "Path must not null")
		resp.InvalidID(w)

		return
	}
	var input entity.UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Err(err).
			Msg("Invalid JSON in update task request")
		// sendResponse(w, "Invalid JSON", false, nil, http.StatusBadRequest)
		resp.InvalidJSON(w)
		return
	}

	err = validation.ValidateStruct(input)
	if err != nil {
		msg := "Validation failed"
		log.Warn().
			Str("request_id", reqID).
			Str("validation_error", err.Error()).
			Msg(msg)

		// sendResponse(w, err.Error(), false, nil, http.StatusBadRequest)
		resp.ValidationError(w, err.Error(), nil)
		return
	}

	newTask := entity.Task{
		Content:   input.Content,
		Completed: input.Completed,
		Timestamp: time.Now().Format(time.RFC3339),
		Priority:  input.Priority,
	}

	task, err := handler.service.UpdateTask(ctx, uint16(intId), newTask)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg("Task not found for update")

		// sendResponse(w, "Task not found", false, nil, http.StatusNotFound)
		resp.TaskNotFound(w, intId)
		return
	}

	task, err = handler.service.GetSingleTask(ctx, uint16(intId))
	if err != nil {
		msg := "Task not found"

		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg(msg)

		// sendResponse(w, msg, false, nil, http.StatusNotFound)
		resp.TaskNotFound(w, intId)

		return
	}

	log.Info().
		Str("request_id", reqID).
		Int("task_id", intId).
		Msg("Task updated successfully")

	// sendResponse(w, "Success", true, updatedTask, http.StatusOK)
	resp.SendSuccessResponse(w, "", task, http.StatusOK)
}

func (handler *DBHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		msg := "Path must not null"
		log.Warn().
			Str("request_id", reqID).
			Str("id", idStr).
			Msg(msg)
		// sendResponse(w, "Path must not null", false, nil, http.StatusBadRequest)
		resp.InternalError(w, msg)
		return
	}

	intId, err := strconv.Atoi(idStr)
	if err != nil {
		msg := "Invalid ID format"
		log.Warn().
			Str("request_id", reqID).
			Str("id", idStr).
			Msg(msg)

		// sendResponse(w, msg, false, nil, http.StatusBadRequest)
		resp.InvalidID(w)
		return
	}

	err = handler.service.DeleteTask(ctx, uint16(intId))
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg("Task not found for deletion")

		// sendResponse(w, "Task not found", false, nil, http.StatusNotFound)
		resp.TaskNotFound(w, intId)

		return
	}

	log.Info().
		Str("request_id", reqID).
		Int("task_id", intId).
		Msg("Task deleted successfully")

	w.WriteHeader(http.StatusNoContent)
}
