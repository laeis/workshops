package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
)

type TaskHandler struct {
	service TaskService
}

type TaskService interface {
	Fetch(ctx context.Context, filters filters.TaskQueryBuilder) (entities.Tasks, error)
	Get(ctx context.Context, id int) (*entities.Task, error)
	Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error)
	Create(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Delete(ctx context.Context, id int) (bool, error)
}

func NewTaskHandler(t TaskService) TaskHandler {
	return TaskHandler{
		service: t,
	}
}

func (t TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	if atoiErr != nil {
		renderErrorResponse(r.Context(), w, atoiErr.Error(), atoiErr)
		return
	}
	task, err := t.service.Get(r.Context(), id)
	if err != nil {
		renderErrorResponse(r.Context(), w, err.Error(), err)
		return
	}
	renderResponse(w, task, http.StatusOK)
}

func (t TaskHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	queryFilters := filters.TaskFilter{}
	queryFilters.Fill(params)
	tasks, err := t.service.Fetch(r.Context(), &queryFilters)
	if err != nil {
		renderErrorResponse(r.Context(), w, err.Error(), err)
		return
	}
	renderResponse(w, tasks, http.StatusOK)
}

func (t TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	newTask := entities.Task{}
	decodeErr := json.NewDecoder(r.Body).Decode(&newTask)
	if decodeErr != nil {
		fmt.Println(decodeErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := t.service.Create(r.Context(), &newTask)
	fmt.Println(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	renderResponse(w, task, http.StatusOK)
}

func (t TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("%v", r.Context())
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	if atoiErr != nil {
		prepareErr := appError.WrapErrorf(atoiErr, appError.ErrorCodeInvalidArgument, "Wrong Id parameter")
		renderErrorResponse(r.Context(), w, "Wrong ID parameter", prepareErr)
		return
	}
	newTask := entities.Task{}
	decodeErr := json.NewDecoder(r.Body).Decode(&newTask)
	if decodeErr != nil {
		renderErrorResponse(r.Context(), w, "Cant decode task data", decodeErr)
		return
	}
	task, err := t.service.Update(r.Context(), id, &newTask)
	if err != nil {
		renderErrorResponse(r.Context(), w, "Task didnt update", err)
		return
	}
	renderResponse(w, task, http.StatusOK)
}

func (t TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, atoiErr := strconv.Atoi(params["id"])
	if atoiErr != nil {
		prepareErr := appError.WrapErrorf(atoiErr, appError.ErrorCodeInvalidArgument, "Wrong Id parameter")
		renderErrorResponse(r.Context(), w, "Wrong ID parameter", prepareErr)
	}
	success, err := t.service.Delete(r.Context(), id)
	if err != nil || !success {
		renderErrorResponse(r.Context(), w, "Task didnt delete", err)
		return
	}
	renderResponse(w, success, http.StatusOK)
}
