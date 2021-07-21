package http

import (
	"github.com/pkg/errors"
	"net/http"
)

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
		next.ServeHTTP(w, r)
	})
}
