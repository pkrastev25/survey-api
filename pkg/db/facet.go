package db

import "go.mongodb.org/mongo-driver/bson"

type FacetBuilder struct {
	conditions []bson.M
}

func NewFacetBuilder() FacetBuilder {
	return FacetBuilder{}
}

func (builder FacetBuilder) Build() []bson.M {
	return builder.conditions
}

func (builder FacetBuilder) Match(condition map[string]interface{}) FacetBuilder {
	builder.conditions = append(builder.conditions, bson.M{operationMatch: condition})
	return builder
}

func (builder FacetBuilder) Sort(condition map[string]interface{}) FacetBuilder {
	builder.conditions = append(builder.conditions, bson.M{operationSort: condition})
	return builder
}

func (builder FacetBuilder) Limit(limit int) FacetBuilder {
	builder.conditions = append(builder.conditions, bson.M{operationLimit: limit})
	return builder
}
