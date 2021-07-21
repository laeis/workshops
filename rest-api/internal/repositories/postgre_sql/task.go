//go:generate mockgen -source taskRepository.go -destination mock/taskRepository_mock.go -package mock
package postgre_sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errWrapper "github.com/pkg/errors"
	"log"
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
	userSql := "SELECT id, title, description, start_date, category  FROM tasks WHERE id = $1"
	if err := r.Connection.QueryRow(userSql, id).Scan(&task.Id, &task.Title, &task.Description, &task.Date, &task.Category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			message := "Task with Id %d not exists"
			log.Println(message, err)
			return nil, fmt.Errorf("Task with Id %d: %w ", id, appError.NotFound)
			//return nil, errWrapper.Wrapf(err, message, id)
		}
		log.Print("Failed to execute query: ", err)
		return nil, err
	}
	return &task, nil
}

func (r *Task) Store(ctx context.Context, task *entities.Task) (*entities.Task, error) {
	query := "INSERT INTO tasks(title, description, category, start_date, user_id) VALUES ($1, $2, $3, $4, $5) returning id"
	var lastInsertId int64
	err := r.Connection.QueryRow(query, task.Title, task.Description, task.Category, task.Date, "6ba7b814-9dad-11d1-80b4-00c04fd430c8").Scan(&lastInsertId)
	if err != nil {
		log.Printf("Error %s when inserting row into tasks table", err)
		return nil, err
	}
	task.Id = int(lastInsertId)
	return task, nil
}

func (r *Task) Update(ctx context.Context, id int, task *entities.Task) (*entities.Task, error) {
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
