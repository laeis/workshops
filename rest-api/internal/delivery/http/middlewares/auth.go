package middlewares

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/metadata"
	"net/http"
	"workshops/rest-api/internal/config"
	rest "workshops/rest-api/internal/delivery/http"
	"workshops/rest-api/internal/delivery/http/response"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/metrics"
	"workshops/rest-api/internal/services"
)

func AuthAdapter(service rest.AuthService, jwtWrapper services.SecurityToken) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientToken, err := jwtWrapper.ParseToken(r)
			if err != nil {
				authError := errors.Wrapf(appError.NotAuthorized, "Wrong casts to claim %w", err)
				response.RenderError(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}
			claims, err := jwtWrapper.ValidateToken(clientToken)
			if err != nil {
				authError := errors.Wrapf(appError.NotAuthorized, "Wrong casts to claim %w", err)
				response.RenderError(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}

			c, ok := claims.(*entities.JwtClaim)
			if !ok {
				authError := errors.Wrap(appError.NotAuthorized, "Wrong casts to claim")
				response.RenderError(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}
			var user *entities.User
			if user, err = service.FindByToken(r.Context(), clientToken); err != nil || user.Email != c.Email {
				authError := errors.Wrapf(appError.NotAuthorized, "Token not exists: %w", err)
				response.RenderError(r.Context(), w, appError.NotAuthorized.Error(), authError)
				return
			}
			ctx := context.WithValue(r.Context(), config.CtxAuthId, user.Id)
			ctx = context.WithValue(ctx, config.CtxToken, clientToken)
			ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
				"Authorization": "Bearer " + clientToken,
			}))

			metrics.RequestUserCounter.With(prometheus.Labels{
				"user_email": user.Email,
				"request":    r.URL.String(),
				"method":     r.Method,
			}).Inc()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
