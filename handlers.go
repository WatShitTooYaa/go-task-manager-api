package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type TaskHandler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *TaskHandler {
	return &TaskHandler{storage: storage}
}

func (handler *TaskHandler) listTask(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	tasks, err := handler.storage.Load()
	if err != nil {
		msg := "Failed to load tasks"

		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)

		// sendResponse(w, msg, false, nil, http.StatusInternalServerError)
		// sendErrorResponse(w, )
		InternalError(w, err.Error())
		return
	}
	// respondSuccess(w, tasks, http.StatusOK)
	msg := "Tasks loaded successfully"
	log.Debug().
		Str("request_id", reqID).
		Int("count", len(tasks)).
		Msg(msg)
	// sendResponse(w, "Success", true, tasks, http.StatusOK)
	sendSuccessResponse(w, msg, tasks, http.StatusOK)
}

func (handler *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	var input CreateTaskRequest
	// input := CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Err(err).
			Msg("Invalid JSON in create task request")
			// sendResponse(w, err.Error(), false, nil, http.StatusInternalServerError)
		InvalidJSON(w)
		return
	}

	if err := ValidateStruct(input); err != nil {
		log.Warn().
			Str("request_id", reqID).
			Str("validation_error", err.Error()).
			Msg("Validation failed")
		// sendResponse(w, err.Error(), false, nil, http.StatusBadRequest)
		ValidationError(w, err.Error(), nil)
		return
	}

	task, err := handler.storage.AddTask(input.Content, input.Priority)
	if err != nil {
		msg := "Failed to create task"
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)

			// sendResponse(w, "Invalid JSON", false, nil, http.StatusInternalServerError)
		InternalError(w, msg)
		return
	}

	log.Info().
		Str("request_id", reqID).
		Uint16("task_id", task.Id).
		Str("content", task.Content).
		Msg("Task created successfully")

		// sendResponse(w, "Task created", true, task, http.StatusCreated) // 201
	sendSuccessResponse(w, "Task created", task, http.StatusCreated)
}

// path
func (handler *TaskHandler) getSingleTask(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	idStr := chi.URLParam(r, "id")
	if idStr == "" {

		// respondError(w, "Path must not null", http.StatusBadRequest)
		// sendResponse(w, "Path must not null", false, nil, http.StatusBadRequest)
		BadRequest(w, "Path must not null")
		return
	}

	intId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Str("id", idStr).
			Msg("Invalid ID format")

		// sendResponse(w, "Path must be int", false, nil, http.StatusBadRequest)
		InvalidID(w)
		return
	}

	task, err := handler.storage.GetByID(uint16(intId))
	if err != nil {
		msg := "Task not found"
		// fmt.Println(msg)

		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg(msg)

		TaskNotFound(w, intId)
		// sendResponse(w, msg, false, nil, http.StatusNotFound)

		return
	}

	log.Debug().
		Str("request_id", reqID).
		Int("task_id", intId).
		Msg("Task retrieved successfully")

	// sendResponse(w, "Success", true, task, http.StatusOK)
	sendSuccessResponse(w, "", task, http.StatusOK)
}

func (handler *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	id := chi.URLParam(r, "id")
	if id == "" {
		// sendResponse(w, "Path must not null", false, nil, http.StatusBadRequest)
		BadRequest(w, "Path must not null")

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
		InvalidID(w)

		return
	}
	var input UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Err(err).
			Msg("Invalid JSON in update task request")
		// sendResponse(w, "Invalid JSON", false, nil, http.StatusBadRequest)
		InvalidJSON(w)
		return
	}

	err = ValidateStruct(input)
	if err != nil {
		msg := "Validation failed"
		log.Warn().
			Str("request_id", reqID).
			Str("validation_error", err.Error()).
			Msg(msg)

		// sendResponse(w, err.Error(), false, nil, http.StatusBadRequest)
		ValidationError(w, err.Error(), nil)
		return
	}

	err = handler.storage.UpdateTask(uint16(intId), input.Content, input.Priority, input.Completed)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg("Task not found for update")

		// sendResponse(w, "Task not found", false, nil, http.StatusNotFound)
		TaskNotFound(w, intId)
		return
	}

	updatedTask, err := handler.storage.GetByID(uint16(intId))
	if err != nil {
		msg := "Task not found"

		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg(msg)

		// sendResponse(w, msg, false, nil, http.StatusNotFound)
		TaskNotFound(w, intId)

		return
	}

	log.Info().
		Str("request_id", reqID).
		Int("task_id", intId).
		Msg("Task updated successfully")

	// sendResponse(w, "Success", true, updatedTask, http.StatusOK)
	sendSuccessResponse(w, "", updatedTask, http.StatusOK)
}

func (handler *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		msg := "Path must not null"
		log.Warn().
			Str("request_id", reqID).
			Str("id", idStr).
			Msg(msg)
		// sendResponse(w, "Path must not null", false, nil, http.StatusBadRequest)
		InternalError(w, msg)
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
		InvalidID(w)
		return
	}

	err = handler.storage.DeleteTask(uint16(intId))
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Int("task_id", intId).
			Msg("Task not found for deletion")

		// sendResponse(w, "Task not found", false, nil, http.StatusNotFound)
		TaskNotFound(w, intId)

		return
	}

	log.Info().
		Str("request_id", reqID).
		Int("task_id", intId).
		Msg("Task deleted successfully")

	w.WriteHeader(http.StatusNoContent)
}
