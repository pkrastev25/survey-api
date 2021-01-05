package pipeline

import (
	"errors"
	dbmodel "survey-api/pkg/db/model"
	"survey-api/pkg/poll/model"
	paginationmodel "survey-api/pkg/poll/pagination/model"

	"go.mongodb.org/mongo-driver/bson"
)

type Builder struct {
	pipeline []bson.M
}

func New() Builder {
	return Builder{
		pipeline: []bson.M{},
	}
}

func (builder Builder) Build() []bson.M {
	return builder.pipeline
}

func (builder Builder) Match(property string, value interface{}) Builder {
	builder.pipeline = append(builder.pipeline, bson.M{string(dbmodel.Match): bson.M{property: value}})
	return builder
}

func (builder Builder) LookUp(from string, localField string, foreignField string, as string) Builder {
	lookUp := bson.M{string(dbmodel.LookUp): bson.M{"from": from, "localField": localField, "foreignField": foreignField, "as": as}}
	builder.pipeline = append(builder.pipeline, lookUp)
	return builder
}

func (builder Builder) Pagination(query paginationmodel.Query) (Builder, error) {
	textSearch := builder.toTextSearch(query.Search())
	if len(textSearch) > 0 {
		builder.pipeline = append(builder.pipeline, textSearch)
	}

	paginate := query.Paginate()
	reversedPaginate, err := builder.reversePaginate(paginate)
	if err != nil {
		return builder, err
	}

	facetStages := bson.M{}
	queryDB, err := builder.toQueryDB(paginate, query.Sort(), query.Limit()+1)
	if err != nil {
		return builder, err
	}

	facetStages[string(paginate.Direction())] = queryDB
	reversedSort, err := builder.reverseSort(query.Sort())
	if err != nil {
		return builder, err
	}

	reversedQueryDB, err := builder.toQueryDB(reversedPaginate, reversedSort, 1)
	if err != nil {
		return builder, err
	}

	facetStages[string(reversedPaginate.Direction())] = reversedQueryDB
	builder.pipeline = append(builder.pipeline, bson.M{string(dbmodel.Facet): facetStages})
	return builder, nil
}

func (builder Builder) ParsePagination(query paginationmodel.Query, paginationPipelineResult []map[string][]model.Poll) ([]model.Poll, map[string]paginationmodel.Query, error) {
	var resultForClient []model.Poll
	paginationQueries := make(map[string]paginationmodel.Query)
	if len(paginationPipelineResult) <= 0 {
		return nil, nil, errors.New("")
	}

	paginationResult := paginationPipelineResult[0]

	paginate := query.Paginate()
	queryResult := paginationResult[string(paginate.Direction())]
	if len(queryResult) == query.Limit()+1 {
		paginationQueries[string(paginate.Direction())] = query.ClonePaginate(builder.generatePaginate(paginate, paginate.Operation(), queryResult[len(queryResult)-1]))
		resultForClient = queryResult[:len(queryResult)-1]
	} else {
		resultForClient = queryResult
	}

	reversePaginate, err := builder.reversePaginate(query.Paginate())
	if err != nil {
		return nil, nil, err
	}

	reverseQueryResult := paginationResult[string(reversePaginate.Direction())]
	if len(reverseQueryResult) > 0 {
		generatedPaginate := builder.generatePaginate(reversePaginate, reversePaginate.Operation(), reverseQueryResult[0])
		generatedSort := builder.generateSort(query.Sort(), generatedPaginate.Direction())
		paginationQueries[string(reversePaginate.Direction())] = query.ClonePaginate(generatedPaginate).CloneSort(generatedSort)
	}

	return resultForClient, paginationQueries, nil
}
