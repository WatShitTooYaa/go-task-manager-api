package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/auth"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	resp "github.com/WatShitTooYaa/go-task-manager-api/internal/response"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/service"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/validation"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s: s}
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	// h.s.LoginService(r.Context())
	var input entity.UserParam
	// input := CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Err(err).
			Msg("Invalid JSON in login request")
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

	user, err := h.s.LoginService(ctx, input)
	if err != nil {
		msg := "Failed to login"
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)

		resp.InternalError(w, msg)
		return
	}

	token, err := auth.GenerateToken(user.Id, user.Password)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate token, err : %s", err.Error())
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)
		resp.InternalError(w, "Failed to generate token")
		return
	}

	log.Info().
		Str("request_id", reqID).
		Uint16("user_id", user.Id).
		Str("username", user.Username).
		Msg("login success")

	resp.SendSuccessResponse(w, "Login successfully", map[string]any{
		"token": token,
	}, http.StatusOK)
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)
	// h.s.LoginService(r.Context())
	var input entity.UserParam
	// input := CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Warn().
			Str("request_id", reqID).
			Err(err).
			Msg("Invalid JSON in register request")
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

	// inputTask := entity.UserParam{

	// }

	user, err := h.s.RegisterService(ctx, input)
	if err != nil {
		msg := fmt.Sprintf("Failed to register. error : %s", err.Error())
		log.Error().
			Str("request_id", reqID).
			Err(err).
			Msg(msg)

		resp.InternalError(w, msg)
		return
	}

	log.Info().
		Str("request_id", reqID).
		Uint16("user_id", user.Id).
		Str("username", user.Username).
		Msg("register success")

	resp.SendSuccessResponse(w, "Register successfully", user, http.StatusCreated)
}
