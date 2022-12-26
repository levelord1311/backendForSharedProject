package storage

import (
	"fmt"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/api/filter"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/api/sort"
)

var _ QueryOptions = &Options{}

var allowedFilters = map[string]string{
	"estate_type": "string",
	"rooms":       "int",
	"district":    "string",
	"price":       "int",
	"created_at":  "date",
	"floor":       "int",
}

func FilterDataType(fltr string) (string, bool) {
	dType, ok := allowedFilters[fltr]
	return dType, ok
}

type Options struct {
	sortField string
	sortOrder string
	fo        map[string][]FilterOption
}

type FilterOption struct {
	Operator string
	Value    []string
	Type     string
}

func NewOptions(so *sort.Options, fo *filter.Options) *Options {

	fltrs := make(map[string][]FilterOption, 0)

	for k, values := range fo.Fields {
		f := FilterOption{}
		for _, v := range values {
			f.Operator = v.Operator
			f.Value = v.Values
			f.Type = v.Type
			fltrs[k] = append(fltrs[k], f)
		}
	}

	return &Options{
		sortField: so.Field,
		sortOrder: so.Order,
		fo:        fltrs,
	}
}

func (o *Options) GetOrderBy() string {
	return fmt.Sprintf("%s %s", o.sortField, o.sortOrder)
}

func (o *Options) GetFilters() map[string][]FilterOption {
	return o.fo
}
