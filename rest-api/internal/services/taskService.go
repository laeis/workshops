package services

import (
	"context"
	"workshops/rest-api/internal/entities"
	"workshops/rest-api/internal/filters"
)

type TaskRepository interface {
	Fetch(ctx context.Context, filters filters.TaskQueryBuilder) (entities.Tasks, error)
	Get(ctx context.Context, id int) (*entities.Task, error)
	Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error)
	Store(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(r TaskRepository) TaskService {
	return TaskService{
		repo: r,
	}
}

func (ts TaskService) Get(ctx context.Context, id int) (*entities.Task, error) {
	return ts.repo.Get(ctx, id)
}
func (ts TaskService) Fetch(ctx context.Context, filters filters.TaskQueryBuilder) (entities.Tasks, error) {
	return ts.repo.Fetch(ctx, filters)
}
func (ts TaskService) Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error) {
	return ts.repo.Update(ctx, id, task)
}
func (ts TaskService) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	return ts.repo.Store(ctx, task)
}
func (ts TaskService) Delete(ctx context.Context, id int) (bool, error) {
	return ts.repo.Delete(ctx, id)
}
