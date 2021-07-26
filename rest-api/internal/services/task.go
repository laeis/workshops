//go:generate mockgen -source taskService.go -destination mock/taskService_mock.go -package mock
package services

import (
	"context"
	"workshops/rest-api/internal/entities"
	"workshops/rest-api/internal/filters"
	"workshops/rest-api/internal/repositories/postgre_sql/builders"
)

type TaskQueryBuilder interface {
	BuildCategoryQuery(string) string
	BuildPeriodQuery(string) string
	BuildOrderQuery(string) string
	BuildOwnerQuery(string, string) string
	QueryArg() []interface{}
}

type TaskRepository interface {
	Fetch(ctx context.Context, filters TaskQueryBuilder) (entities.Tasks, error)
	Get(ctx context.Context, id int) (*entities.Task, error)
	Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error)
	Store(ctx context.Context, task *entities.Task) (*entities.Task, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type TaskService struct {
	repo TaskRepository
	user UserRepository
}

func NewTask(r TaskRepository, ur UserRepository) *TaskService {
	return &TaskService{
		repo: r,
		user: ur,
	}
}

func (ts *TaskService) Get(ctx context.Context, id int) (*entities.Task, error) {
	task, err := ts.repo.Get(ctx, id)
	if err != nil {
		return task, err
	}
	timezone, err := ts.user.GetAuthTimezone(ctx)
	if err != nil {
		return nil, err
	}
	task.InTimezone(timezone)
	return task, err
}
func (ts *TaskService) Fetch(ctx context.Context, filters *filters.TaskFilter) (entities.Tasks, error) {
	builder := builders.NewTask(filters)
	tasks, err := ts.repo.Fetch(ctx, &builder)
	timezone, err := ts.user.GetAuthTimezone(ctx)
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		tasks[i].InTimezone(timezone)
	}
	return tasks, err
}
func (ts *TaskService) Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error) {
	return ts.repo.Update(ctx, id, task)
}
func (ts *TaskService) Create(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	return ts.repo.Store(ctx, task)
}
func (ts *TaskService) Delete(ctx context.Context, id int) (bool, error) {
	return ts.repo.Delete(ctx, id)
}
