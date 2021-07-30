package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"workshops/rest-api/internal/delivery/http/response"
	"workshops/rest-api/internal/entities"
	appErrors "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/metrics"
	"workshops/rest-api/internal/validators"
)

type UserHandler struct {
	service UserService
	log     *zap.Logger
}

type UserService interface {
	Get(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, id string, user *entities.User) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
}

func NewUser(s UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: s,
		log:     logger,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	userTemplate := entities.EmptyUser()

	decodeErr := json.NewDecoder(r.Body).Decode(&userTemplate)
	if decodeErr != nil {
		response.RenderError(r.Context(), w, appErrors.BadRequest.Error(), errors.Wrapf(appErrors.BadRequest, "Couldnt decode request: %w", decodeErr))
		return
	}
	validator := validators.UserValidator(&userTemplate)
	if err := validator.Validate(); err != nil {
		response.RenderError(r.Context(), w, err.Error(), errors.Wrapf(appErrors.BadRequest, "Validation error: %w", err))
		return
	}

	user, err := h.service.Create(r.Context(), &userTemplate)
	if err != nil {
		response.RenderError(r.Context(), w, "", err)
		return
	}
	metrics.UserCnt.Inc()
	response.Render(w, user, http.StatusOK)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		message := "Id is required parameter"
		response.RenderError(r.Context(), w, message, errors.Wrap(appErrors.BadRequest, message))
		return
	}

	userTemplate := entities.EmptyUser()
	decodeErr := json.NewDecoder(r.Body).Decode(&userTemplate)
	if decodeErr != nil {
		response.RenderError(r.Context(), w, appErrors.BadRequest.Error(), errors.Wrapf(appErrors.BadRequest, "Couldnt decode request: %w", decodeErr))
		return
	}

	validator := validators.UserValidator(&userTemplate)
	if err := validator.Validate("timezone"); err != nil {
		response.RenderError(r.Context(), w, err.Error(), errors.Wrapf(appErrors.BadRequest, "Validation error: %w", err))
		return
	}

	user, err := h.service.Update(r.Context(), id, &userTemplate)
	if err != nil {
		response.RenderError(r.Context(), w, "", err)
		return
	}

	response.Render(w, user, http.StatusOK)
}
