package response

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	appError "workshops/rest-api/internal/errors"
)

//Response represents a response containing an error message.
type responsePayload struct {
	Error   string      `json:"error,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func RenderError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	errorMsg := msg
	if errorMsg == "" {
		errorMsg = err.Error()
	}
	resp := responsePayload{Error: errorMsg}
	status := http.StatusInternalServerError
	switch true {
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, appError.NotFound):
		status = http.StatusNotFound
	case errors.Is(err, appError.BadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, appError.AccessForbidden):
		status = http.StatusForbidden
	case errors.Is(err, appError.NotAuthorized):
		status = http.StatusUnauthorized
	}
	Render(w, resp, status)
}

func Render(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	data := res
	if _, ok := data.(responsePayload); !ok {
		data = responsePayload{
			Payload: res,
		}
	}

	content, err := json.Marshal(data)
	if err != nil {
		// XXX Do something with the error ;)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// XXX Do something with the error ;)
	}
}
