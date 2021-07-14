package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"workshops/rest-api/internal/entities"
	appError "workshops/rest-api/internal/errors"
	"workshops/rest-api/internal/filters"
)

type Task struct {
	Connection *sql.DB
}

func NewTask(c *sql.DB) *Task {
	return &Task{Connection: c}
}

func (r *Task) Fetch(ctx context.Context, filter filters.TaskQueryBuilder) (entities.Tasks, error) {
	var count int
	err := r.Connection.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	query := "SELECT id, title, description, start_date, category FROM tasks WHERE 0 = 0 "
	query = filter.BuildCategoryQuery(query)
	query = filter.BuildPeriodQuery(query)
	query = filter.BuildOrderQuery(query)

	stmt, err := r.Connection.Prepare(query)
	if err != nil {
		log.Print("Failed to prepare query: ", err)
		return nil, err
	}

	rows, err := stmt.Query(filter.QueryArg()...)
	if err != nil {
		log.Print("Failed to execute query: ", err)
		return nil, err
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
	err := r.Connection.QueryRow(userSql, id).Scan(&task.Id, &task.Title, &task.Description, &task.Date, &task.Category)
	if err != nil {
		switch true {
		case errors.Is(err, sql.ErrNoRows):
			message := "Task with Id %d not exists"
			log.Print(message, err)
			return nil, appError.WrapErrorf(err, appError.ErrorCodeNotFound, message, id)
		default:
			log.Print("Failed to execute query: ", err)
			return nil, err
		}

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
	query, err := r.Connection.Prepare("update tasks set title=$1, description=$2, category=$3, start_date=$4 where id=$5")
	if err != nil {
		log.Printf("Error %s when prepare query for update task", err)
		return nil, err
	}

	_, err = query.Exec(task.Title, task.Description, task.Category, task.Date, id)
	if err != nil {
		log.Printf("Error %s when update task %d", err, id)
		return nil, err
	}
	uTask, err := r.Get(ctx, id)
	if err != nil {
		log.Printf("Error %s when get updatable task %d", err, id)
		return nil, err
	}
	return uTask, nil
}

func (r *Task) Delete(ctx context.Context, id int) (bool, error) {
	stmt, err := r.Connection.PrepareContext(ctx, "delete from tasks where id=$1")
	if err != nil {
		log.Printf("Error %s when prepare query for delete task %d", err, id)
		return false, err
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		log.Printf("Error %s when delete task %d", err, id)
		return false, err
	}
	return true, nil
}
