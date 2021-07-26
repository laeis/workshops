package builders

type Storage interface {
	QueryArg() []interface{}
}

type ArgStorage struct {
	args []interface{}
}

func (as *ArgStorage) QueryArg() []interface{} {
	return as.args
}

func (as *ArgStorage) AddQueryArg(value ...interface{}) {
	as.args = append(as.args, value...)
}

type QueryPlaceholder struct {
	placeholder int
}

func (qp *QueryPlaceholder) NextPlaceholder() int {
	qp.placeholder++
	return qp.placeholder
}
