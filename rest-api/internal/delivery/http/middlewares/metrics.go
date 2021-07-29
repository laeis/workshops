package middlewares

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"workshops/rest-api/internal/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(w http.ResponseWriter) {
			metrics.RequestCounter.With(prometheus.Labels{
				"request": r.URL.String(),
				"status":  w.Header().Get("status"),
				"method":  r.Method,
			}).Inc()
		}(w)

		next.ServeHTTP(w, r)
	})
}
