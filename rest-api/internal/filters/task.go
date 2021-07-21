package filters

type TaskFilterValidator interface {
	ValidateParameter(name string, value []string) string
}

type TaskFilter struct {
	Category string
	Period   string
	Order    string `json:"order,omitempty"`
	OrderBy  string `json:"order_by,omitempty"`
}

func ValidatedTaskFilter(v TaskFilterValidator, params map[string][]string) (f TaskFilter) {
	f.Category = v.ValidateParameter("category", params["category"])
	f.Period = v.ValidateParameter("period", params["period"])
	f.Order = v.ValidateParameter("order", params["order"])
	f.OrderBy = v.ValidateParameter("order_by", params["order_by"])
	return
}
