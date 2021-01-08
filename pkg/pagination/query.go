package pagination

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	defaultLimit = 20
	minLimit     = 1
	maxLimit     = 100
)

type Query struct {
	paginate Paginate
	sort     Sort
	limit    int
}

func NewQuery() Query {
	return Query{
		paginate: NewPaginate(),
		sort:     NewSort(),
		limit:    defaultLimit,
	}
}

func (query Query) ClonePaginate(paginate Paginate) Query {
	query.paginate = paginate
	return query
}

func (query Query) CloneSort(sort Sort) Query {
	query.sort = sort
	return query
}

func (query Query) Paginate() Paginate {
	return query.paginate
}

func (query Query) Sort() Sort {
	return query.sort
}

func (query Query) Limit() int {
	return query.limit
}

func (query *Query) SetPaginate(paginate Paginate) {
	query.paginate = paginate
}

func (query *Query) SetSort(sort Sort) {
	query.sort = sort
}

func (query *Query) SetLimit(limit int) {
	if limit < minLimit {
		query.limit = defaultLimit
	} else if limit > maxLimit {
		query.limit = maxLimit
	} else {
		query.limit = limit
	}
}

func (query Query) Validate() error {
	return validation.ValidateStruct(&query,
		validation.Field(&query.paginate),
		validation.Field(&query.sort),
		validation.Field(&query.limit, validation.Min(minLimit), validation.Max(maxLimit)),
	)
}
