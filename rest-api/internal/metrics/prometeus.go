package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TaskCnt = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rest_task_total",
		Help: "Current task counter .",
	})

	UserCnt = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rest_user_total",
		Help: "Current user counter.",
	})

	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rest_http_request_total",
			Help: "Number of request to rest service.",
		},
		[]string{"request", "method", "status"},
	)

	RequestUserCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rest_http_request_user_total",
			Help: "Number of auth request to rest service per user.",
		},
		[]string{"user_email", "request", "method"},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter, RequestUserCounter, TaskCnt, UserCnt)
}
