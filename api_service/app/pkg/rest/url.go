package rest

import (
	"fmt"
	"strings"
)

type FilterOptions struct {
	Field    string
	Operator string
	Values   []string
}

// ToStringWF provides filtering options like in:1,3,4 or neq:4 or eq:1 or =123
func (f *FilterOptions) ToStringWF() string {
	return fmt.Sprintf("%s%s", f.Operator, strings.Join(f.Values, ","))
}
