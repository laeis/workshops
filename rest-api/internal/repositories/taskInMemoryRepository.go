package repositories

import (
	"context"
	"errors"
	"time"
	"workshops/rest-api/internal/entities"
	"workshops/rest-api/internal/filters"
)

type InMemoryTask struct {
	Data entities.Tasks
}

func NewInMemoryTask() *InMemoryTask {
	return &InMemoryTask{Data: []entities.Task{
		{
			Id:          1,
			Title:       "Task Title",
			Description: "Some Task Description",
			Date:        time.Now(),
		},
	}}
}

func (r *InMemoryTask) Fetch(ctx context.Context, filters filters.TaskQueryBuilder) (entities.Tasks, error) {
	return r.Data, nil
}

func (r *InMemoryTask) Get(ctx context.Context, id int) (*entities.Task, error) {
	for _, t := range r.Data {
		if t.Id == id {
			return &t, nil
		}
	}
	return nil, errors.New("Task not exists. ")
}
func (r *InMemoryTask) Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error) {
	for i, _ := range r.Data {
		if r.Data[i].Id == id {
			r.Data[i].Title = task.Title
			r.Data[i].Description = task.Description
			r.Data[i].Date = task.Date
			return &r.Data[i], nil
		}
	}
	return nil, errors.New("Task not exists. ")
}

func (r *InMemoryTask) Store(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	r.Data = append(r.Data, *task)
	return &r.Data[len(r.Data)-1], nil
}

func (r *InMemoryTask) Delete(ctx context.Context, id int) (bool, error) {
	for i, _ := range r.Data {
		if r.Data[i].Id == id {
			r.Data = append(r.Data[:i], r.Data[i+1:]...)
			return true, nil
		}
	}
	return false, errors.New("Task not exists. ")
}
