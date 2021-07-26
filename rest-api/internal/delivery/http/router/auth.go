package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func Auth(r *mux.Router, h AuthHandler, authMiddleware mux.MiddlewareFunc) {
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)

	logout := r.PathPrefix("/logout").Subrouter()
	logout.Use(authMiddleware)
	logout.HandleFunc("", h.Logout)
}
