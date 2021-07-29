//go:generate mockgen -source taskHandler.go -destination mock/taskHandler_mock.go -package mock
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"workshops/rest-api/internal/delivery/http/response"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
	"workshops/rest-api/internal/metrics"
	"workshops/rest-api/internal/validators"
)

type TaskHandler struct {
	service TaskService
}

type TaskService interface {
	Fetch(ctx context.Context, filters *filters.TaskFilter) (entities.Tasks, error)
	Get(ctx context.Context, id int) (*entities.Task, error)
	Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error)
	Create(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Delete(ctx context.Context, id int) (bool, error)
}

func NewTask(t TaskService) TaskHandler {
	return TaskHandler{
		service: t,
	}
}

func (t TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])

	if atoiErr != nil {
		//fmt.Println( errors.Is(&appError.NotFound{Err: atoiErr}, &appError.BadRequest{}))
		response.RenderError(r.Context(), w, atoiErr.Error(), appError.BadRequest)
		return
	}

	task, err := t.service.Get(r.Context(), id)
	if err != nil {
		response.RenderError(r.Context(), w, err.Error(), err)
		return
	}

	response.Render(w, task, http.StatusOK)
}

func (t TaskHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	validator := validators.TaskValidator{}
	queryFilters := filters.ValidatedTaskFilter(&validator, r.URL.Query())
	tasks, err := t.service.Fetch(r.Context(), &queryFilters)
	if err != nil {
		response.RenderError(r.Context(), w, err.Error(), err)
		return
	}
	response.Render(w, tasks, http.StatusOK)
}

func (t TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	newTask := entities.Task{}
	decodeErr := json.NewDecoder(r.Body).Decode(&newTask)
	if decodeErr != nil {
		response.RenderError(r.Context(), w, "Wrong data for new task", fmt.Errorf("%q: %w", "Wrong data for new task", appError.BadRequest))
		return
	}

	task, err := t.service.Create(r.Context(), &newTask)
	if err != nil {
		response.RenderError(r.Context(), w, err.Error(), err)
		return
	}
	metrics.TaskCnt.Inc()
	response.Render(w, task, http.StatusOK)
}

func (t TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	if atoiErr != nil {
		message := "Wrong Id parameter"
		response.RenderError(r.Context(), w, message, fmt.Errorf("%q: %w", atoiErr.Error(), appError.BadRequest))
		return
	}

	newTask := entities.Task{}
	decodeErr := json.NewDecoder(r.Body).Decode(&newTask)
	if decodeErr != nil {
		response.RenderError(r.Context(), w, "Cant decode task data", fmt.Errorf("%q: %w", decodeErr.Error(), appError.BadRequest))
		return
	}

	task, err := t.service.Update(r.Context(), id, &newTask)
	if err != nil {
		response.RenderError(r.Context(), w, "Task didnt update", err)
		return
	}

	response.Render(w, task, http.StatusOK)
}

func (t TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	if atoiErr != nil {
		message := "Wrong Id parameter"
		prepareErr := fmt.Errorf("%q: %w", message, appError.BadRequest)
		response.RenderError(r.Context(), w, message, prepareErr)
		return
	}

	success, err := t.service.Delete(r.Context(), id)
	if err != nil || !success {
		response.RenderError(r.Context(), w, "Task didnt delete", err)
		return
	}
	metrics.TaskCnt.Dec()
	response.Render(w, success, http.StatusOK)
}
