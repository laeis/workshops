package http

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/delivery/http/response"
	"workshops/rest-api/internal/entities"
	appErrors "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/validators"
)

type AuthHandler struct {
	service AuthService
}

type AuthService interface {
	Login(ctx context.Context, user entities.User) (*entities.JWT, error)
	Logout(ctx context.Context, email string, token string) (bool, error)
	FindByToken(ctx context.Context, token string) (*entities.User, error)
}

func NewAuth(s AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	userTemplate := entities.EmptyUser()

	decodeErr := json.NewDecoder(r.Body).Decode(&userTemplate)
	if decodeErr != nil {
		response.RenderError(r.Context(), w, appErrors.BadRequest.Error(), errors.Wrapf(appErrors.BadRequest, "Couldnt decode request: %w", decodeErr))
		return
	}
	validator := validators.UserValidator(&userTemplate)
	if err := validator.Validate("email", "password"); err != nil {
		response.RenderError(r.Context(), w, err.Error(), errors.Wrapf(appErrors.BadRequest, "Validation error: %w", err))
		return
	}

	token, err := h.service.Login(r.Context(), userTemplate)
	if err != nil {
		response.RenderError(r.Context(), w, err.Error(), errors.Wrapf(appErrors.BadRequest, "Wrong credentional %w", err))
		return
	}

	response.Render(w, token, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, idOk := r.Context().Value(config.CtxAuthId).(string)
	token, tokenOk := r.Context().Value(config.CtxToken).(string)
	if !tokenOk && !idOk {
		message := "Required parameters are missing"
		response.RenderError(r.Context(), w, message, errors.Wrap(appErrors.BadRequest, message))
		return
	}

	success, err := h.service.Logout(r.Context(), id, token)
	if err != nil {
		response.RenderError(r.Context(), w, err.Error(), errors.Wrapf(appErrors.BadRequest, "Some error hapen %w", err))
		return
	}

	response.Render(w, success, http.StatusOK)
}
