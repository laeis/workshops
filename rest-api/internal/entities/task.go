package entities

import (
	"time"
)

type Tasks []Task

type Task struct {
	Id          int
	Title       string
	Description string
	Category    string
	UserId      string
	Date        time.Time
}

func (t *Task) InTimezone(timezone string) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return
	}
	t.Date = t.Date.In(location)
}

func (t *Task) UpdateIfExists(task *Task) {
	if task.Title != "" {
		t.Title = task.Title
	}
	if task.Description != "" {
		t.Description = task.Description
	}
	if task.Category != "" {
		t.Category = task.Category
	}
	if !task.Date.IsZero() && task.Date != time.Unix(0, 0).UTC() {
		t.Date = task.Date
	}
}
