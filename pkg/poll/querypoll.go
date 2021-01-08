package poll

import (
	"survey-api/pkg/pagination"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type QueryPoll struct {
	base   pagination.Query
	search string
}

func NewQueryPoll() QueryPoll {
	return QueryPoll{
		base: pagination.NewQuery(),
	}
}

func (queryPoll QueryPoll) ClonePaginate(paginate pagination.Paginate) QueryPoll {
	queryPoll.base = queryPoll.base.ClonePaginate(paginate)
	return queryPoll
}

func (queryPoll QueryPoll) CloneSort(sort pagination.Sort) QueryPoll {
	queryPoll.base = queryPoll.base.CloneSort(sort)
	return queryPoll
}

func (queryPoll QueryPoll) Search() string {
	return queryPoll.search
}

func (queryPoll *QueryPoll) SetBase(base pagination.Query) {
	queryPoll.base = base
}

func (queryPoll *QueryPoll) SetSearch(search string) {
	queryPoll.search = search
}

func (queryPoll QueryPoll) Validate() error {
	err := validation.ValidateStruct(&queryPoll.base)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&queryPoll,
		validation.Field(&queryPoll.search, is.Alphanumeric),
	)
}
