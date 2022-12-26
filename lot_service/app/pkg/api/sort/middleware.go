package sort

import (
	"context"
	"net/http"
	"strings"
)

const (
	ASC               = "ASC"
	DESC              = "DESC"
	OptionsContextKey = "sort_options"
	DefSort           = "created_at"
	DefOrder          = "DESC"
)

type Options struct {
	Field, Order string
}

// Middleware parses query for sorting parameters 'sort_by' and 'sort_order'. If absent, uses DefSort and DefOrder.
func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = DefSort
		}

		if sortOrder == "" {
			sortOrder = DefOrder
		} else {
			upperSortOrder := strings.ToUpper(sortOrder)
			if upperSortOrder != ASC && upperSortOrder != DESC {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("bad sorting order"))
				// TODO w.Write - error for your API
				return
			}
		}

		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}
		ctx := context.WithValue(r.Context(), OptionsContextKey, options)
		r = r.WithContext(ctx)

		h(w, r)
	}
}
