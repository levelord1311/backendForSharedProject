package filter

var allowedOperators = map[string]string{
	"eq":  "=",
	"neq": "!=",
	"lt":  "<",
	"lte": "<=",
	"gt":  ">",
	"gte": ">=",
}

type Options struct {
	Fields map[string][]Field
}
type Field struct {
	Operator string
	Values   []string
	Type     string
}

func NewOptions(fields map[string][]Field) *Options {
	return &Options{
		Fields: fields,
	}
}

func OperatorIsAllowed(o string) (string, bool) {
	op, ok := allowedOperators[o]
	return op, ok
}
