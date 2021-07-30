//go:generate mockgen -source taskHandler.go -destination mock/taskHandler_mock.go -package mock
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/delivery/http/response"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
	"workshops/rest-api/internal/metrics"
	"workshops/rest-api/internal/validators"
)

type TaskHandler struct {
	service TaskService
	log     *zap.Logger
}

type TaskService interface {
	Fetch(ctx context.Context, filters *filters.TaskFilter) (entities.Tasks, error)
	Get(ctx context.Context, id int) (*entities.Task, error)
	Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error)
	Create(ctx context.Context, task *entities.Task, userId string) (*entities.Task, error)
	Delete(ctx context.Context, id int) (bool, error)
}

func NewTask(t TaskService, logger *zap.Logger) TaskHandler {
	return TaskHandler{
		service: t,
		log:     logger,
	}
}

func (t TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	logger := t.log.With(zap.String("handler", "GetTask"), zap.String("route", r.URL.RequestURI()))

	if atoiErr != nil {
		wErr := errors.Wrap(appError.BadRequest, "Could`nt decode id")
		logger.Error(wErr.Error(), zap.String("id", params["id"]))
		response.RenderError(r.Context(), w, atoiErr.Error(), wErr)
		return
	}

	task, err := t.service.Get(r.Context(), id)

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

func (t TaskHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	logger := t.log.With(zap.String("handler", "FetchTask"), zap.String("route", r.URL.RequestURI()))
	validator := validators.TaskValidator{}
	queryFilters := filters.ValidatedTaskFilter(&validator, r.URL.Query())
	tasks, err := t.service.Fetch(r.Context(), &queryFilters)

	if err != nil {
		wErr := errors.Wrap(err, "Could`nt get tasks")
		logger.Error(wErr.Error(), zap.String("query", r.URL.Query().Encode()))
		response.RenderError(r.Context(), w, err.Error(), wErr)
		return
	}

	response.Render(w, tasks, http.StatusOK)
}

func (t TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	logger := t.log.With(zap.String("handler", "CreateTask"), zap.String("route", r.URL.RequestURI()))
	newTask := entities.Task{}
	decodeErr := json.NewDecoder(r.Body).Decode(&newTask)

	if decodeErr != nil {
		message := "Wrong data for new task"
		wErr := errors.Wrap(appError.BadRequest, message)
		logger.Error(wErr.Error())
		response.RenderError(r.Context(), w, message, wErr)
		return
	}

	authId, ok := r.Context().Value(config.CtxAuthId).(string)

	if !ok || authId == "" {
		logger.Error(appError.NotAuthorized.Error())
		response.RenderError(r.Context(), w, appError.NotAuthorized.Error(), appError.NotAuthorized)
		return
	}

	task, err := t.service.Create(r.Context(), &newTask, authId)

	if err != nil {
		message := "Task not created"
		wErr := errors.Wrap(err, message)
		logger.Error(wErr.Error())
		response.RenderError(r.Context(), w, message, wErr)
		return
	}

	metrics.TaskCnt.Inc()
	response.Render(w, task, http.StatusOK)
}

func (t TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	logger := t.log.With(zap.String("handler", "UpdateTask"), zap.String("route", r.URL.RequestURI()))
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])

	if atoiErr != nil {
		message := "Wrong Id parameter"
		logger.Error(message, zap.Any("id", params["id"]))
		response.RenderError(r.Context(), w, message, fmt.Errorf("%q: %w", atoiErr.Error(), appError.BadRequest))
		return
	}

	newTask := entities.Task{}
	decodeErr := json.NewDecoder(r.Body).Decode(&newTask)

	if decodeErr != nil {
		message := "Cant decode task data"
		wErr := errors.Wrapf(appError.BadRequest, "%s : %w", message, decodeErr)
		logger.Error(message)
		response.RenderError(r.Context(), w, message, wErr)
		return
	}

	task, err := t.service.Update(r.Context(), id, &newTask)

	if err != nil {
		message := "Task didnt update"
		logger.Error(message, zap.Int("id", id))
		response.RenderError(r.Context(), w, "Task didnt update", err)
		return
	}

	response.Render(w, task, http.StatusOK)
}

func (t TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	logger := t.log.With(zap.String("handler", "DeleteTask"), zap.String("route", r.URL.RequestURI()))
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])

	if atoiErr != nil {
		message := "Wrong Id parameter"
		prepareErr := fmt.Errorf("%q: %w", message, appError.BadRequest)
		logger.Error(message, zap.Int("id", id))
		response.RenderError(r.Context(), w, message, prepareErr)
		return
	}

	success, err := t.service.Delete(r.Context(), id)

	if err != nil || !success {
		message := "Task didnt delete"
		logger.Error(message, zap.Int("id", id))
		response.RenderError(r.Context(), w, message, err)
		return
	}

	metrics.TaskCnt.Dec()

	response.Render(w, success, http.StatusOK)
}
