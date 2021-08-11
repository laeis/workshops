//go:generate mockgen -source taskHandler.go -destination mock/taskHandler_mock.go -package mock
package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"workshops/rest-api/internal/delivery/http/response"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
	"workshops/rest-api/internal/validators"
	"workshops/rest-api/pb"
)

type TaskGRPCHandler struct {
	client pb.TaskServiceClient
	log    *zap.Logger
}

func NewTaskGRPC(c pb.TaskServiceClient, logger *zap.Logger) TaskGRPCHandler {
	return TaskGRPCHandler{
		client: c,
		log:    logger,
	}
}

func (t TaskGRPCHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	logger := t.log.With(zap.String("handler", "GetTask"), zap.String("route", r.URL.RequestURI()))

	if atoiErr != nil {
		wErr := errors.Wrap(appError.BadRequest, "Could`nt decode id")
		logger.Error(wErr.Error(), zap.String("id", params["id"]))
		response.RenderError(r.Context(), w, atoiErr.Error(), wErr)
		return
	}

	task, err := t.client.Get(r.Context(), &pb.TaskRequest{
		Id: uint32(id),
	})

	if err != nil {
		wErr := errors.Wrap(err, "Could`nt get task by id")
		if err == appError.NotFound {
			wErr = err
		}

		logger.Error(wErr.Error(), zap.String("id", params["id"]))
		response.RenderError(r.Context(), w, wErr.Error(), wErr)
		return
	}

	response.Render(w, task, http.StatusOK)
}

func (t TaskGRPCHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	logger := t.log.With(zap.String("handler", "FetchTask"), zap.String("route", r.URL.RequestURI()))
	validator := validators.TaskValidator{}
	queryFilters := filters.ValidatedTaskFilter(&validator, r.URL.Query())
	bq, err := json.Marshal(queryFilters)
	fetchRequest := pb.FetchTaskRequest{}
	err = json.Unmarshal(bq, &fetchRequest)

	tasks, err := t.client.Fetch(r.Context(), &fetchRequest)

	if err != nil {
		wErr := errors.Wrap(err, "Could`nt get tasks")
		logger.Error(wErr.Error(), zap.String("query", r.URL.Query().Encode()))
		response.RenderError(r.Context(), w, err.Error(), wErr)
		return
	}

	response.Render(w, tasks, http.StatusOK)
}
