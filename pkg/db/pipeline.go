package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

type PipelineBuilder struct {
	pipeline []bson.M
}

func NewPipelineBuilder() PipelineBuilder {
	return PipelineBuilder{}
}

func (builder PipelineBuilder) Build() []bson.M {
	return builder.pipeline
}

func (builder PipelineBuilder) FacetStage(stages map[string]FacetBuilder) PipelineBuilder {
	facet := bson.M{}
	for name, query := range stages {
		facet[name] = query.Build()
	}

	builder.pipeline = append(builder.pipeline, bson.M{operationFacet: facet})
	return builder
}

func (builder PipelineBuilder) MatchStage(property string, value interface{}) PipelineBuilder {
	builder.pipeline = append(builder.pipeline, bson.M{operationMatch: bson.M{property: value}})
	return builder
}

func (builder PipelineBuilder) LookUpStage(from string, localField string, foreignField string, as string) PipelineBuilder {
	builder.pipeline = append(builder.pipeline, bson.M{operationLookUp: bson.M{"from": from, "localField": localField, "foreignField": foreignField, "as": as}})
	return builder
}

func (builder PipelineBuilder) TextSearchStage(search string) PipelineBuilder {
	builder.pipeline = append(builder.pipeline, bson.M{operationText: bson.M{operationSearch: search}})
	return builder
}

func (builder PipelineBuilder) SetStage(property string, value interface{}) PipelineBuilder {
	builder.pipeline = append(builder.pipeline, bson.M{operationSet: bson.M{property: value}})
	return builder
}
