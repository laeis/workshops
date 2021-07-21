package validators

var validCategories = stringList{NOTE, EVENT}

var validPeriods = stringList{DAY, WEEK, MONTH, YEAR}

var validOrder = stringList{ASC, DESC}

var validOrderField = stringList{"id", "title", "start_date"}

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
		if a == v {
			return true
		}
	}
	return false
}

type TaskValidator struct {
}

func (tv *TaskValidator) ValidateParameter(name string, value []string) string {
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
