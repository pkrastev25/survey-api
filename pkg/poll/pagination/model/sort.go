package model

import (
	"errors"
	dbmodel "survey-api/pkg/db/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "des"
)

type SortOrder string

type Sort struct {
	order    SortOrder
	property dbmodel.Property
}

func NewSort() Sort {
	return Sort{
		order:    Ascending,
		property: dbmodel.Created,
	}
}

func NewSortOrder(input string) (SortOrder, error) {
	if input == string(Ascending) {
		return Ascending, nil
	}

	if input == string(Descending) {
		return Descending, nil
	}

	return "", errors.New("")
}

func NewSortProperty(input string) (dbmodel.Property, error) {
	if input == string(dbmodel.Created) {
		return dbmodel.Created, nil
	}

	return "", errors.New("")
}

func NewSortPaginate(paginate Paginate) (Sort, error) {
	value := Ascending
	if paginate.Direction() == Prev {
		value = Descending
	}

	return NewSortAll(string(paginate.Property()), string(value))
}

func NewSortAll(sourceProperty string, sourceOrder string) (Sort, error) {
	var sort Sort
	property, err := NewSortProperty(sourceProperty)
	if err != nil {
		return sort, err
	}

	order, err := NewSortOrder(sourceOrder)
	if err != nil {
		return sort, err
	}

	return Sort{property: property, order: order}, nil
}

func (sort Sort) Property() dbmodel.Property {
	return sort.property
}

func (sort Sort) Order() SortOrder {
	return sort.order
}

func (sort Sort) Validate() error {
	return validation.ValidateStruct(&sort,
		validation.Field(&sort.property, validation.In(dbmodel.Created)),
		validation.Field(&sort.order, validation.In(Ascending, Descending)),
	)
}

func (s Sort) CloneOrder(order SortOrder) Sort {
	s.order = order
	return s
}
