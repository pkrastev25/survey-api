package mapper

import (
	"errors"
	dbmodel "survey-api/pkg/db/model"
	"survey-api/pkg/poll/pagination/model"

	"go.mongodb.org/mongo-driver/bson"
)

func ToSearchDB(search string) bson.M {
	if len(search) <= 0 {
		return nil
	}

	return bson.M{string(dbmodel.Text): bson.M{string(dbmodel.Search): search}}
}

func ToPaginateDB(paginate model.Paginate) (bson.M, error) {
	value, err := paginate.FormatValue()
	if err != nil {
		return nil, err
	}

	return bson.M{string(paginate.Property()): bson.M{string(paginate.Operation()): value}}, nil
}

func ToSortDB(sort model.Sort) (bson.M, error) {
	if sort.Order() == model.Ascending {
		return bson.M{string(sort.Property()): 1}, nil
	}

	if sort.Order() == model.Descending {
		return bson.M{string(sort.Property()): -1}, nil
	}

	return nil, errors.New("")
}
