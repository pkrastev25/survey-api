package pagination

import (
	"survey-api/pkg/db"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	Ascending  = "asc"
	Descending = "des"
)

type Sort struct {
	order    string
	property string
}

func NewSort() Sort {
	return Sort{
		order:    Ascending,
		property: db.PropertyCreated,
	}
}

func NewSortPaginate(paginate Paginate) Sort {
	order := Ascending
	if paginate.Direction() == Prev {
		order = Descending
	}

	return NewSortAll(paginate.Property(), order)
}

func NewSortAll(property string, order string) Sort {
	return Sort{property: property, order: order}
}

func (s Sort) CloneReverseOrder() Sort {
	if s.order == Descending {
		return s.CloneOrder(Ascending)
	}

	return s.CloneOrder(Descending)
}

func (s Sort) CloneOrder(order string) Sort {
	s.order = order
	return s
}

func (sort Sort) Property() string {
	return sort.property
}

func (sort Sort) Order() string {
	return sort.order
}

func (sort Sort) Validate() error {
	return validation.ValidateStruct(&sort,
		validation.Field(&sort.property, validation.In(db.PropertyCreated)),
		validation.Field(&sort.order, validation.In(Ascending, Descending)),
	)
}
