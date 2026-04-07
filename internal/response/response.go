package response

import (
	"encoding/json"
	"net/http"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/utils"
)

type ErrorCode string

const (
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidJSON  ErrorCode = "INVALID_JSON"
	ErrCodeInvalidID    ErrorCode = "INVALID_ID"
	ErrCodeTaskNotFound ErrorCode = "TASK_NOT_FOUND"
	ErrCodeStorageError ErrorCode = "STORAGE_ERROR"
)

type ErrorDetail struct {
	Code    ErrorCode      `json:"code"`
	Message string         `json:"message"`
	Details utils.JsonType `json:"details,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// sendErrorResponse sends structured error response
func sendErrorResponse(w http.ResponseWriter, code ErrorCode, message string, statusCode int, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// sendSuccessResponse sends structured success response
func SendSuccessResponse(w http.ResponseWriter, message string, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Helper functions for common errors
func BadRequest(w http.ResponseWriter, message string) {
	sendErrorResponse(w, ErrCodeBadRequest, message, http.StatusBadRequest, nil)
}

func ValidationError(w http.ResponseWriter, message string, details map[string]interface{}) {
	sendErrorResponse(w, ErrCodeValidation, message, http.StatusBadRequest, details)
}

func InvalidJSON(w http.ResponseWriter) {
	sendErrorResponse(w, ErrCodeInvalidJSON, "Invalid JSON format", http.StatusBadRequest, nil)
}

func InvalidID(w http.ResponseWriter) {
	sendErrorResponse(w, ErrCodeInvalidID, "Invalid ID format", http.StatusBadRequest, nil)
}

func TaskNotFound(w http.ResponseWriter, id int) {
	sendErrorResponse(w, ErrCodeTaskNotFound, "Task not found", http.StatusNotFound, map[string]interface{}{
		"task_id": id,
	})
}

func InternalError(w http.ResponseWriter, message string) {
	sendErrorResponse(w, ErrCodeInternal, message, http.StatusInternalServerError, nil)
}
