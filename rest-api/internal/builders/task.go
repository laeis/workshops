package builders

import (
	"fmt"
	"time"
	"workshops/rest-api/internal/filters"
)

type Task struct {
	args        []interface{}
	placeholder int
	filters     *filters.TaskFilter
}

func (tb *Task) NextPlaceholder() int {
	tb.placeholder++
	return tb.placeholder
}

func NewTask(filters *filters.TaskFilter) Task {
	return Task{filters: filters}
}

func (tb *Task) BuildCategoryQuery(query string) string {
	if tb.filters.Category == "" {
		return query
	}
	tb.args = append(tb.args, tb.filters.Category)
	return fmt.Sprintf("%s AND category=$%d ", query, tb.NextPlaceholder())
}

func (tb *Task) BuildPeriodQuery(query string) string {
	if tb.filters.Period == "" {
		return query
	}
	beginPlaceholder := tb.NextPlaceholder()
	endPlaceholder := tb.NextPlaceholder()
	periodQuery := fmt.Sprintf("%s AND (start_date >= $%d  AND start_date <= $%d) ", query, beginPlaceholder, endPlaceholder)
	begin, end := createPeriod(tb.filters.Period)
	tb.args = append(tb.args, begin, end)
	return periodQuery
}

func (tb *Task) BuildOrderQuery(query string) string {
	return fmt.Sprintf("%s ORDER BY %s %s", query, tb.filters.OrderBy, tb.filters.Order)
}
func (tb *Task) QueryArg() []interface{} {
	return tb.args
}

func createPeriod(period string) (time.Time, time.Time) {
	y, m, d := time.Now().Date()
	var begin, end time.Time
	switch period {
	case "day":
		begin = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		end = time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
	case "week":
		currentDay := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		weekday := int(currentDay.Weekday())
		weekStartDayInt := int(time.Monday)
		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
		begin = currentDay.AddDate(0, 0, -weekday)
		end = begin.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case "month":
		begin = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
		end = begin.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case "year":
		begin = time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
		end = begin.AddDate(1, 0, 0).Add(-time.Nanosecond)
	}
	return begin, end
}
