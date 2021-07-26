package http

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/services"
)

type Adapter func(http.Handler) http.Handler

func AuthMiddlewareAdapter(service AuthService, jwtWrapper services.SecurityToken) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientToken, err := jwtWrapper.ParseToken(r)
			if err != nil {
				authError := errors.Wrapf(appError.NotAuthorized, "Wrong casts to claim %w", err)
				renderErrorResponse(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}
			claims, err := jwtWrapper.ValidateToken(clientToken)
			if err != nil {
				authError := errors.Wrapf(appError.NotAuthorized, "Wrong casts to claim %w", err)
				renderErrorResponse(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}

			c, ok := claims.(*entities.JwtClaim)
			if !ok {
				authError := errors.Wrap(appError.NotAuthorized, "Wrong casts to claim")
				renderErrorResponse(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}
			var user *entities.User
			if user, err = service.FindByToken(r.Context(), clientToken); err != nil || user.Email != c.Email {
				authError := errors.Wrapf(appError.NotAuthorized, "Token not exists: %w", err)
				renderErrorResponse(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}
			ctx := context.WithValue(r.Context(), config.CtxAuthId, user.Id)
			ctx = context.WithValue(ctx, config.CtxToken, clientToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec != nil {
				var err error
				switch t := rec.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				renderErrorResponse(r.Context(), w, err.Error(), err)
			}
		}()
		fmt.Println("recover")
		next.ServeHTTP(w, r)
	})
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
