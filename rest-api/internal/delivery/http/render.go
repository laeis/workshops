package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	appError "workshops/rest-api/internal/errors"
)

// ErrorResponse represents a response containing an error message.
type Response struct {
	Error   string      `json:"error,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func renderErrorResponse(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	errorMsg := msg
	if errorMsg == "" {
		errorMsg = err.Error()
	}
	resp := Response{Error: errorMsg}
	status := http.StatusInternalServerError
	switch true {
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, appError.NotFound):
		status = http.StatusNotFound
	case errors.Is(err, appError.BadRequest):
		status = http.StatusBadRequest
	}
	renderResponse(w, resp, status)
}

func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	data := res
	if _, ok := data.(Response); !ok {
		data = Response{
			Payload: res,
		}
	}

	content, err := json.Marshal(data)
	if err != nil {
		// XXX Do something with the error ;)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// XXX Do something with the error ;)
	}
}
