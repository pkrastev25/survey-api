package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

type QueryBuilder struct {
	query bson.M
}

func NewQueryBuilder() QueryBuilder {
	return QueryBuilder{
		query: bson.M{},
	}
}

func (builder QueryBuilder) Build() bson.M {
	return builder.query
}

func (builder QueryBuilder) Equal(property string, value interface{}) QueryBuilder {
	builder.query[property] = value
	return builder
}

func (builder QueryBuilder) Set(property string, value interface{}) QueryBuilder {
	builder.query[operationSet] = bson.M{property: value}
	return builder
}

func (builder QueryBuilder) NotIn(property string, value []interface{}) QueryBuilder {
	builder.query[property] = bson.M{operationNotIn: value}
	return builder
}

func (builder QueryBuilder) AddToSet(property string, value interface{}) QueryBuilder {
	builder.query[operationAddToSet] = bson.M{property: value}
	return builder
}

func (builder QueryBuilder) Increment(property string, increment int) QueryBuilder {
	builder.query[operationIncrement] = bson.M{property: increment}
	return builder
}
