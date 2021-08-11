package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type TaskHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
}

func TaskGRPC(r *mux.Router, tc TaskHandler, authMiddleware mux.MiddlewareFunc) {
	r.Use(authMiddleware)
	r.HandleFunc("/{id}", tc.Get).Methods(http.MethodGet)
	r.HandleFunc("/", tc.Fetch).Methods(http.MethodGet)
}
