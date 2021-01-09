package pagination

import (
	"survey-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
)

type PaginationMapper struct {
}

func NewPaginationMapper() PaginationMapper {
	return PaginationMapper{}
}

func (mapper PaginationMapper) PaginateToDb(paginate Paginate) (bson.M, error) {
	value, err := paginate.FormatValue()
	if err != nil {
		return nil, err
	}

	return mapper.paginateToDb(paginate.Property(), paginate.Operation(), value), nil
}

func (mapper PaginationMapper) ReversePaginateToDb(paginate Paginate) (bson.M, error) {
	value, err := paginate.FormatValue()
	if err != nil {
		return nil, err
	}

	reverseOperation := db.LessThan
	if db.LessThanOrEqual == paginate.Property() {
		reverseOperation = db.GreaterThan
	}

	return mapper.paginateToDb(paginate.Property(), reverseOperation, value), nil
}

func (mapper PaginationMapper) paginateToDb(property string, operation string, value interface{}) bson.M {
	return bson.M{property: bson.M{operation: value}}
}

func (mapper PaginationMapper) SortToDb(sort Sort) bson.M {
	order := 1
	if sort.Order() == Descending {
		order = -1
	}

	return bson.M{sort.Property(): order}
}
