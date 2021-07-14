package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Controller interface {
	Get(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewRouter(tc Controller) http.Handler {
	r := mux.NewRouter()
	TaskHandler := tc

	r.HandleFunc("/tasks/{id}", TaskHandler.Get).Methods(http.MethodGet)
	r.HandleFunc("/tasks", TaskHandler.Fetch).Methods(http.MethodGet)
	r.HandleFunc("/tasks", TaskHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id}", TaskHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/tasks/{id}", TaskHandler.Delete).Methods(http.MethodDelete)

	return r
}
