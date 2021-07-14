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
	Date        time.Time
}
