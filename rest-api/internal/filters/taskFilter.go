package filters

import (
	"fmt"
	"time"
)

const (
	NOTE  = "note"
	EVENT = "event"
	ASC   = "asc"
	DESC  = "desc"
)

const (
	DAY   = "day"
	WEEK  = "week"
	MONTH = "month"
	YEAR  = "year"
)

type stringList []string

func (c stringList) Contains(v string) bool {
	for _, a := range c {
		fmt.Println(a, v)
		if a == v {
			return true
		}
	}
	return false
}

var validCategories = stringList{NOTE, EVENT}

var validPeriods = stringList{DAY, WEEK, MONTH, YEAR}

var validOrder = stringList{ASC, DESC}

var validOrderField = stringList{"id", "title", "start_date"}

type TaskQueryBuilder interface {
	BuildCategoryQuery(string) string
	BuildPeriodQuery(string) string
	BuildOrderQuery(string) string
	QueryArg() []interface{}
}

type TaskBuilder struct {
	args        []interface{}
	placeholder int
}

func (tb *TaskBuilder) NextPlaceholder() int {
	tb.placeholder++
	return tb.placeholder
}

type TaskFilter struct {
	Category string
	Period   string
	Order    string `json:"order,omitempty"`
	OrderBy  string `json:"order_by,omitempty"`
	builder  TaskBuilder
}

func (f *TaskFilter) Fill(params map[string][]string) {
	f.Category = f.validateParameters("category", params["category"])
	f.Period = f.validateParameters("period", params["period"])
	f.Order = f.validateParameters("order", params["order"])
	f.OrderBy = f.validateParameters("order_by", params["order_by"])
}

func (f *TaskFilter) validateParameters(name string, value []string) string {
	switch name {
	case "category":
		if len(value) > 0 && validCategories.Contains(value[0]) {
			return value[0]
		}
	case "period":
		if len(value) > 0 && validPeriods.Contains(value[0]) {
			return value[0]
		}
	case "order":
		if len(value) > 0 && validOrder.Contains(value[0]) {
			return value[0]
		} else {
			return ASC
		}
	case "order_by":
		if len(value) > 0 && validOrderField.Contains(value[0]) {
			return value[0]
		} else {
			return "id"
		}
	}
	return ""
}

func (f *TaskFilter) BuildCategoryQuery(query string) string {
	if f.Category == "" {
		return query
	}
	f.builder.args = append(f.builder.args, f.Category)
	return fmt.Sprintf("%s AND category=$%d ", query, f.builder.NextPlaceholder())
}

func (f *TaskFilter) BuildPeriodQuery(query string) string {
	if f.Period == "" {
		return query
	}
	beginPlaceholder := f.builder.NextPlaceholder()
	endPlaceholder := f.builder.NextPlaceholder()
	periodQuery := fmt.Sprintf("%s AND (start_date >= $%d  AND start_date <= $%d) ", query, beginPlaceholder, endPlaceholder)
	begin, end := createPeriod(f.Period)
	f.builder.args = append(f.builder.args, begin, end)
	return periodQuery
}

func (f *TaskFilter) BuildOrderQuery(query string) string {
	return fmt.Sprintf("%s ORDER BY %s %s", query, f.OrderBy, f.Order)
}
func (f *TaskFilter) QueryArg() []interface{} {
	return f.builder.args
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
