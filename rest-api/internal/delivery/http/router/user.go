package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

func User(r *mux.Router, h UserHandler, authMiddleware mux.MiddlewareFunc) {
	r.Use(authMiddleware)
	r.HandleFunc("/", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id}", h.Update).Methods(http.MethodPut)
}
