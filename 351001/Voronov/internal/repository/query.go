package repository

type Pagination struct {
	Page     int
	PageSize int
}

type Filter struct {
	Field string
	Value interface{}
}

type Sort struct {
	Field     string
	Direction string
}

type QueryOptions struct {
	Pagination *Pagination
	Filters    []Filter
	Sort       *Sort
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Pagination: &Pagination{Page: 1, PageSize: 10},
		Filters:    []Filter{},
		Sort:       &Sort{Field: "id", Direction: "ASC"},
	}
}

func WithPagination(page, pageSize int) func(*QueryOptions) {
	return func(o *QueryOptions) {
		if page > 0 {
			o.Pagination.Page = page
		}
		if pageSize > 0 {
			o.Pagination.PageSize = pageSize
		}
	}
}

func WithFilter(field string, value interface{}) func(*QueryOptions) {
	return func(o *QueryOptions) {
		o.Filters = append(o.Filters, Filter{Field: field, Value: value})
	}
}

func WithSort(field, direction string) func(*QueryOptions) {
	return func(o *QueryOptions) {
		o.Sort = &Sort{Field: field, Direction: direction}
	}
}
