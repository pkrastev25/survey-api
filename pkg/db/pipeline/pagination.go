package pipeline

import (
	"errors"
	dbmodel "survey-api/pkg/db/model"
	"survey-api/pkg/dtime"
	"survey-api/pkg/poll/model"
	paginationmodel "survey-api/pkg/poll/pagination/model"
	"survey-api/pkg/poll/pagination/model/mapper"

	"go.mongodb.org/mongo-driver/bson"
)

func (builder Builder) reversePaginate(paginate paginationmodel.Paginate) (paginationmodel.Paginate, error) {
	var reversePaginate paginationmodel.Paginate
	if paginate.Operation() == dbmodel.GreaterThan {
		return paginate.CloneOperation(dbmodel.LessThanOrEqual).CloneDirection(paginationmodel.Prev), nil
	}

	if paginate.Operation() == dbmodel.GreaterThanOrEqual {
		return paginate.CloneOperation(dbmodel.LessThan).CloneDirection(paginationmodel.Prev), nil
	}

	if paginate.Operation() == dbmodel.LessThan {
		return paginate.CloneOperation(dbmodel.GreaterThanOrEqual).CloneDirection(paginationmodel.Next), nil
	}

	if paginate.Operation() == dbmodel.LessThanOrEqual {
		return paginate.CloneOperation(dbmodel.GreaterThan).CloneDirection(paginationmodel.Next), nil
	}

	return reversePaginate, errors.New("")
}

func (builder Builder) reverseSort(sort paginationmodel.Sort) (paginationmodel.Sort, error) {
	var reverseSort paginationmodel.Sort
	if sort.Order() == paginationmodel.Ascending {
		return sort.CloneOrder(paginationmodel.Descending), nil
	}

	if sort.Order() == paginationmodel.Descending {
		return sort.CloneOrder(paginationmodel.Ascending), nil
	}

	return reverseSort, errors.New("")
}

func (builder Builder) toQueryDB(paginate paginationmodel.Paginate, sort paginationmodel.Sort, limit int) ([]bson.M, error) {
	var queriesDB []bson.M
	matchDB, err := builder.toMatchDB(paginate)
	if err != nil {
		return nil, err
	}

	if len(matchDB) > 0 {
		queriesDB = append(queriesDB, matchDB)
	}

	sortDB, err := mapper.ToSortDB(sort)
	if err != nil {
		return nil, err
	}

	if len(sortDB) > 0 {
		queriesDB = append(queriesDB, bson.M{string(dbmodel.Sort): sortDB})
	}

	if limit <= 0 {
		return nil, errors.New("")
	}

	queriesDB = append(queriesDB, builder.toLimitDB(limit))
	return queriesDB, nil
}

func (builder Builder) toTextSearch(search string) bson.M {
	searchDB := mapper.ToSearchDB(search)
	if len(searchDB) <= 0 {
		return nil
	}

	return bson.M{string(dbmodel.Match): searchDB}
}

func (builder Builder) toMatchDB(paginate paginationmodel.Paginate) (bson.M, error) {
	paginateDB, err := mapper.ToPaginateDB(paginate)
	if err != nil {
		return nil, err
	}

	if len(paginateDB) > 0 {
		return bson.M{string(dbmodel.Match): paginateDB}, nil
	}

	return nil, nil
}

func (builder Builder) toLimitDB(limit int) bson.M {
	return bson.M{string(dbmodel.Limit): limit}
}

func (builder Builder) generateSort(sort paginationmodel.Sort, direction paginationmodel.PaginateDirection) paginationmodel.Sort {
	order := paginationmodel.Ascending
	if direction == paginationmodel.Prev {
		order = paginationmodel.Descending
	}

	return sort.CloneOrder(order)
}

func (builder Builder) generatePaginate(sourcePaginate paginationmodel.Paginate, sourceOperation dbmodel.Operation, poll model.Poll) paginationmodel.Paginate {
	operation := dbmodel.GreaterThanOrEqual
	if sourceOperation == dbmodel.LessThan || sourceOperation == dbmodel.LessThanOrEqual {
		operation = dbmodel.LessThanOrEqual
	}

	paginate := sourcePaginate.CloneOperation(operation)
	switch paginate.Property() {
	case dbmodel.Created:
		paginate = paginate.CloneValue(dtime.DateTimeToISO(poll.Created))
	}

	return paginate
}
