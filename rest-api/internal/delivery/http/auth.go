package http

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/delivery/http/response"
	"workshops/rest-api/internal/entities"
	appErrors "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/validators"
)

type AuthHandler struct {
	service AuthService
	log     *zap.Logger
}

type AuthService interface {
	Login(ctx context.Context, user entities.User) (*entities.JWT, error)
	Logout(ctx context.Context, email string, token string) (bool, error)
	FindByToken(ctx context.Context, token string) (*entities.User, error)
}

func NewAuth(s AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: s,
		log:     logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := h.log.With(zap.String("handler", "Login"), zap.String("route", r.URL.RequestURI()))

	userTemplate := entities.EmptyUser()
	decodeErr := json.NewDecoder(r.Body).Decode(&userTemplate)

	if decodeErr != nil {
		wErr := errors.Wrapf(appErrors.BadRequest, "Couldnt decode request: %w", decodeErr)
		logger.Error(wErr.Error())
		response.RenderError(
			r.Context(),
			w,
			appErrors.BadRequest.Error(),
			wErr,
		)
		return
	}

	validator := validators.UserValidator(&userTemplate)

	if err := validator.Validate("email", "password"); err != nil {
		wErr := errors.Wrapf(appErrors.BadRequest, "Validation error: %w", err)
		logger.Error(wErr.Error(), zap.String("email", userTemplate.Email), zap.String("password", userTemplate.Password))
		response.RenderError(
			r.Context(),
			w,
			err.Error(),
			wErr,
		)
		return
	}

	token, err := h.service.Login(r.Context(), userTemplate)

	if err != nil {
		wErr := errors.Wrapf(appErrors.BadRequest, "Wrong credentional %w", err)
		logger.Error(wErr.Error(), zap.String("email", userTemplate.Email))
		response.RenderError(
			r.Context(),
			w,
			err.Error(),
			wErr,
		)
		return
	}

	response.Render(w, token, http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := h.log.With(zap.String("handler", "Logout"), zap.String("route", r.URL.RequestURI()))

	id, idOk := r.Context().Value(config.CtxAuthId).(string)
	token, tokenOk := r.Context().Value(config.CtxToken).(string)

	if !tokenOk && !idOk {
		message := "Required parameters are missing"
		wErr := errors.Wrap(appErrors.BadRequest, message)
		logger.Error(wErr.Error())
		response.RenderError(r.Context(), w, message, wErr)
		return
	}

	success, err := h.service.Logout(r.Context(), id, token)

	if err != nil {
		wErr := errors.Wrapf(appErrors.BadRequest, "Some error hapen %w", err)
		response.RenderError(
			r.Context(),
			w,
			err.Error(),
			wErr,
		)
		return
	}

	response.Render(w, success, http.StatusOK)
}
