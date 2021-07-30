//go:generate mockgen -source taskRepository.go -destination mock/taskRepository_mock.go -package mock
package postgre_sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errWrapper "github.com/pkg/errors"
	"log"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/services"
)

type Task struct {
	Connection *sql.DB
}

func NewTask(c *sql.DB) *Task {
	return &Task{Connection: c}
}

func (r *Task) Fetch(ctx context.Context, taskBuilder services.TaskQueryBuilder) (entities.Tasks, error) {
	var count int
	err := r.Connection.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)

	query := "SELECT id, title, description, start_date, category FROM tasks WHERE 0 = 0 "
	query = taskBuilder.BuildCategoryQuery(query)
	query = taskBuilder.BuildPeriodQuery(query)
	userId := ctx.Value(config.CtxAuthId).(string)
	query = taskBuilder.BuildOwnerQuery(query, userId)
	query = taskBuilder.BuildOrderQuery(query)

	rows, err := r.Connection.Query(query, taskBuilder.QueryArg()...)
	defer rows.Close()

	if err != nil {
		return nil, errWrapper.Wrap(err, "Failed to execute query")
	}

	tasks := make(entities.Tasks, 0, count)

	for rows.Next() {
		var task entities.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Date, &task.Category)

		if err != nil {
			log.Print("Failed to scan task: ", err)
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Task) Get(ctx context.Context, id int) (*entities.Task, error) {
	task := entities.Task{}
	userSql := "SELECT id, title, description, start_date, category, user_id  FROM tasks WHERE id = $1"

	if err := r.Connection.QueryRow(userSql, id).Scan(&task.Id, &task.Title, &task.Description, &task.Date, &task.Category, &task.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errWrapper.Wrapf(appError.NotFound, "Task with Id %d not exists", id)
		}
		return nil, err
	}

	authId := ctx.Value(config.CtxAuthId)

	if task.UserId != authId {
		err := appError.AccessForbidden
		return nil, err
	}

	return &task, nil
}

func (r *Task) Store(ctx context.Context, task *entities.Task, userId string) (*entities.Task, error) {
	query := "INSERT INTO tasks(title, description, category, start_date, user_id) VALUES ($1, $2, $3, $4, $5) returning id"

	var lastInsertId int64
	err := r.Connection.QueryRow(query, task.Title, task.Description, task.Category, task.Date, userId).Scan(&lastInsertId)

	if err != nil {
		return nil, err
	}

	task.Id = int(lastInsertId)
	return task, nil
}

func (r *Task) Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error) {
	if !r.hasRequestPermission(ctx, id) {
		err := appError.AccessForbidden
		log.Printf("Permission denied: %s", err)
		return nil, err
	}

	query := "update tasks set title=$1, description=$2, category=$3, start_date=$4 where id=$5"
	_, err := r.Connection.Exec(query, task.Title, task.Description, task.Category, task.Date, id)

	if err != nil {
		log.Printf("Error %s when update task %d", err, id)
		return nil, err
	}

	uTask, err := r.Get(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Task with Id %d: %w ", id, appError.NotFound)
		}
		log.Print("Failed to execute query: ", err)
		return nil, err
	}

	return uTask, nil
}

func (r *Task) Delete(ctx context.Context, id int) (bool, error) {
	if !r.hasRequestPermission(ctx, id) {
		err := appError.AccessForbidden
		log.Printf("Permission denied: %s", err)
		return false, err
	}
	result, err := r.Connection.ExecContext(ctx, "delete from tasks where id=$1", id)

	if result != nil {
		d, _ := result.RowsAffected()
		if d == 0 {
			return false, fmt.Errorf("Task with Id %d: %w ", id, appError.NotFound)
		}
	}

	if err != nil {
		log.Printf("Error %s when delete task %d", err, id)
		return false, err
	}

	return true, nil
}

func (r *Task) hasRequestPermission(ctx context.Context, taskId int) bool {
	task := entities.Task{}
	authId := ctx.Value(config.CtxAuthId)
	userSql := "SELECT id FROM tasks WHERE id = $1 AND user_id = $2"

	if err := r.Connection.QueryRow(userSql, taskId, authId).Scan(&task.Id); err != nil {
		return false
	}

	return true
}
