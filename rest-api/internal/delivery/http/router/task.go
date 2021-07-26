package router

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

func Task(r *mux.Router, tc Controller, authMiddleware mux.MiddlewareFunc) {
	r.Use(authMiddleware)
	r.HandleFunc("/{id}", tc.Get).Methods(http.MethodGet)
	r.HandleFunc("/", tc.Fetch).Methods(http.MethodGet)
	r.HandleFunc("/", tc.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id}", tc.Update).Methods(http.MethodPut)
	r.HandleFunc("/{id}", tc.Delete).Methods(http.MethodDelete)
}
