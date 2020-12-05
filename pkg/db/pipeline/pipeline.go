package pipeline

import (
	"errors"
	"survey-api/pkg/poll/model"
	"survey-api/pkg/poll/pagination"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	GreaterThan Operation = "$gte"
	LessThan    Operation = "$lt"
	Limit       Operation = "$limit"
	Facet       Operation = "$facet"
	Match       Operation = "$match"
)

type Operation string

type Builder struct {
	pipeline []bson.M
}

func New() *Builder {
	return &Builder{
		pipeline: []bson.M{},
	}
}

func (builder *Builder) Pagination(metadata *pagination.Metadata) *Builder {
	facet := bson.M{string(Facet): metadata.ToDbPagination()}
	builder.pipeline = append(builder.pipeline, facet)

	return builder
}

func CreateLimitCondition(limit int) bson.M {
	return bson.M{string(Limit): limit}
}

func (builder *Builder) ParsePagination(pipelineResult []map[string][]model.Poll) (map[string][]model.Poll, error) {
	if len(pipelineResult) <= 0 {
		return nil, errors.New("")
	}

	return pipelineResult[0], nil
}

func (builder *Builder) Build() interface{} {
	return builder.pipeline
}
