package model

import (
	"errors"
	dbmodel "survey-api/pkg/db/model"
	"survey-api/pkg/dtime"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	Next PaginateDirection = "next"
	Prev PaginateDirection = "prev"
)

type PaginateDirection string

type Paginate struct {
	property  dbmodel.Property
	operation dbmodel.Operation
	value     string
	direction PaginateDirection
}

func NewPaginateProperty(input string) (dbmodel.Property, error) {
	if input == string(dbmodel.Created) {
		return dbmodel.Created, nil
	}

	return "", errors.New("")
}

func NewPaginateOperation(input string) (dbmodel.Operation, error) {
	if input == string(dbmodel.GreaterThanOrEqual) {
		return dbmodel.GreaterThanOrEqual, nil
	}

	if input == string(dbmodel.LessThanOrEqual) {
		return dbmodel.LessThanOrEqual, nil
	}

	return "", errors.New("")
}

func NewPaginate() Paginate {
	return Paginate{
		property:  dbmodel.Created,
		operation: dbmodel.GreaterThanOrEqual,
		value:     dtime.NilTimeISO,
		direction: Next,
	}
}

func NewPaginateAll(sourceProperty string, sourceOperation string, value string) (Paginate, error) {
	var paginate Paginate
	property, err := NewPaginateProperty(sourceProperty)
	if err != nil {
		return paginate, err
	}

	operation, err := NewPaginateOperation(sourceOperation)
	if err != nil {
		return paginate, err
	}

	direction := Next
	if operation == dbmodel.LessThanOrEqual {
		direction = Prev
	}

	return Paginate{
		property:  property,
		operation: operation,
		value:     value,
		direction: direction,
	}, nil
}

func (p Paginate) Property() dbmodel.Property {
	return p.property
}

func (p Paginate) Operation() dbmodel.Operation {
	return p.operation
}

func (p Paginate) Value() string {
	return p.value

}

func (p Paginate) FormatValue() (interface{}, error) {
	switch p.property {
	case dbmodel.Created:
		return dtime.ISOToDateTime(p.value)
	default:
		return nil, errors.New("")
	}
}

func (p Paginate) Direction() PaginateDirection {
	return p.direction
}

func (p Paginate) CloneValue(value string) Paginate {
	p.value = value
	return p
}

func (p Paginate) CloneDirection(direction PaginateDirection) Paginate {
	p.direction = direction
	return p
}

func (p Paginate) CloneOperation(operation dbmodel.Operation) Paginate {
	p.operation = operation
	return p
}

func (p Paginate) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.property, validation.In(dbmodel.Created)),
		validation.Field(&p.operation, validation.In(dbmodel.GreaterThanOrEqual, dbmodel.LessThanOrEqual)),
		validation.Field(&p.value, validation.Date(dtime.ISOFormat)),
	)
}
