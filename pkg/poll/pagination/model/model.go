package model

import (
	"errors"
	"survey-api/pkg/db/pipeline"
	"survey-api/pkg/dtime"
	"survey-api/pkg/poll/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type Query struct {
	Filter []FilterCondition `json:"filter"`
	Sort   []SortCondition   `json:"sort"`
}

type Condition struct {
	Property  string      `json:"property"`
	Operation string      `json:"operation,omitempty"`
	Value     interface{} `json:"value"`
}

type FilterCondition struct {
	Condition
}

type SortCondition struct {
	Condition
}

func DefaultDbCondition() bson.M {
	return bson.M{}
}

func (q Query) Validate() error {
	return validation.ValidateStruct(&q)
}

func (c FilterCondition) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Property, validation.Required, validation.In("created")),
		validation.Field(&c.Operation, validation.In(string(pipeline.GreaterThan), string(pipeline.LessThan))),
		validation.Field(&c.Value, validation.Required, validation.By(validateValue(c))),
	)
}

func validateValue(c FilterCondition) validation.RuleFunc {
	return func(value interface{}) error {
		switch c.Property {
		case "created":
			t, err := time.Parse(dtime.ISOFormat, value.(string))
			if err != nil {
				return err
			}

			c.Value = t
		default:
			return errors.New("")
		}

		return nil
	}
}

func (q *Query) NewFromPoll(poll *model.Poll) *Query {
	var filters []FilterCondition

	if len(q.Filter) == 0 {
		filters = append(filters, FilterCondition{})
	} else {
		for _, condition := range q.Filter {
			filters = append(filters, *condition.NewFromPoll(poll))
		}
	}

	return &Query{
		Filter: filters,
		Sort:   q.Sort,
	}
}

func (fc *FilterCondition) NewFromPoll(poll *model.Poll) *FilterCondition {
	var value interface{}

	switch fc.Property {
	case "created":
		value = poll.Created.Time().Format(dtime.ISOFormat)
	default:
		return nil
	}

	condition := &FilterCondition{}
	condition.Property = fc.Property
	condition.Operation = fc.Operation
	condition.Value = value
	return condition
}

func (q *Query) ToDbQuery() ([]bson.M, []bson.M) {
	var next []bson.M
	var prev []bson.M

	for _, item := range q.Sort {
		sort := item.toDbSort()
		next = append(next, sort)
		prev = append(next, sort)
	}

	for _, item := range q.Filter {
		next = append(next, item.toDbFilterNext())
		prev = append(prev, item.toDbFilterPrev())
	}

	return next, prev
}

func (sc *SortCondition) toDbSort() bson.M {
	return sc.toDbCondition(sc.Operation)
}

func (fc *FilterCondition) toDbFilterNext() bson.M {
	return fc.toDbCondition(fc.Operation)
}

func (fc *FilterCondition) toDbFilterPrev() bson.M {
	prevOperation := ""

	if fc.Operation == string(pipeline.GreaterThan) {
		prevOperation = string(pipeline.LessThan)
	} else if fc.Operation == string(pipeline.LessThan) {
		prevOperation = string(pipeline.GreaterThan)
	}

	return fc.toDbCondition(prevOperation)
}

func (c *Condition) toDbCondition(operation string) bson.M {
	dbCondition := bson.M{}
	inner := c.Value

	if len(operation) > 0 {
		inner = bson.M{
			operation: c.Value,
		}
	}

	dbCondition[string(pipeline.Match)] = bson.M{
		c.Property: inner,
	}

	return dbCondition
}
