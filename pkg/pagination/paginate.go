package pagination

import (
	"errors"
	"survey-api/pkg/db"
	"survey-api/pkg/dtime"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	Next = "next"
	Prev = "prev"
)

type Paginate struct {
	property  string
	operation string
	value     string
	direction string
}

func NewPaginate() Paginate {
	return Paginate{
		property:  db.Created,
		operation: db.GreaterThanOrEqual,
		value:     dtime.NilTimeISO,
		direction: Next,
	}
}

func NewPaginateAll(property string, operation string, value string) Paginate {
	direction := Next
	if operation == db.LessThanOrEqual {
		direction = Prev
	}

	return Paginate{
		property:  property,
		operation: operation,
		value:     value,
		direction: direction,
	}
}

func (paginate Paginate) CloneValue(value string) Paginate {
	paginate.value = value
	return paginate
}

func (paginate Paginate) CloneDirection(direction string) Paginate {
	paginate.direction = direction
	return paginate
}

func (paginate Paginate) CloneReverseDirection() Paginate {
	if paginate.direction == Prev {
		return paginate.CloneDirection(Next)
	}

	return paginate.CloneDirection(Prev)
}

func (paginate Paginate) CloneOperation(operation string) Paginate {
	paginate.operation = operation
	return paginate
}

func (paginate Paginate) CloneReverseOperation() Paginate {
	if paginate.operation == db.GreaterThanOrEqual {
		return paginate.CloneOperation(db.LessThanOrEqual)
	}

	return paginate.CloneOperation(db.GreaterThanOrEqual)
}

func (paginate Paginate) Property() string {
	return paginate.property
}

func (paginate Paginate) Operation() string {
	return paginate.operation
}

func (paginate Paginate) Value() string {
	return paginate.value
}

func (paginate Paginate) FormatValue() (interface{}, error) {
	switch paginate.property {
	case db.Created:
		return dtime.ISOToDateTime(paginate.value)
	default:
		return nil, errors.New("")
	}
}

func (paginate Paginate) Direction() string {
	return paginate.direction
}

func (paginate Paginate) Validate() error {
	return validation.ValidateStruct(&paginate,
		validation.Field(&paginate.property, validation.In(db.Created)),
		validation.Field(&paginate.operation, validation.In(db.GreaterThanOrEqual, db.LessThanOrEqual)),
		validation.Field(&paginate.value, validation.Date(dtime.ISOFormat)),
	)
}
