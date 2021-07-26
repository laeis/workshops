package entities

import (
	"fmt"
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
	fmt.Println(t.Date)
	fmt.Println(t.Date.In(location))
	t.Date = t.Date.In(location)
}
