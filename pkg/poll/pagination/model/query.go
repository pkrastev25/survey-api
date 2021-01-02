package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	defaultLimit = 20
	minLimit     = 1
	maxLimit     = 100
)

type Query struct {
	search   string
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

func (q Query) ClonePaginate(paginate Paginate) Query {
	q.paginate = paginate
	return q
}

func (q Query) CloneSort(sort Sort) Query {
	q.sort = sort
	return q
}

func (q Query) Search() string {
	return q.search
}

func (q Query) Paginate() Paginate {
	return q.paginate
}

func (q Query) Sort() Sort {
	return q.sort
}

func (q Query) Limit() int {
	return q.limit
}

func (q *Query) SetSearch(search string) {
	q.search = search
}

func (q *Query) SetPaginate(paginate Paginate) {
	q.paginate = paginate
}

func (q *Query) SetSort(sort Sort) {
	q.sort = sort
}

func (q *Query) SetLimit(limit int) {
	if limit < minLimit {
		q.limit = defaultLimit
	} else if limit > maxLimit {
		q.limit = maxLimit
	} else {
		q.limit = limit
	}
}

func (q Query) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.search, is.Alphanumeric),
		validation.Field(&q.paginate),
		validation.Field(&q.sort),
		validation.Field(&q.limit, validation.Min(minLimit), validation.Max(maxLimit)),
	)
}
