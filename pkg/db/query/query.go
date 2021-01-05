package query

import (
	"survey-api/pkg/db/model"

	"go.mongodb.org/mongo-driver/bson"
)

type Builder struct {
	query bson.M
}

func New() Builder {
	return Builder{
		query: bson.M{},
	}
}

func NewMap(conditions map[string]interface{}) Builder {
	return Builder{
		query: conditions,
	}
}

func (builder Builder) Build() bson.M {
	return builder.query
}

func (builder Builder) Filter(property string, value interface{}) Builder {
	builder.query[property] = value
	return builder
}

func (builder Builder) Update(property string, value interface{}) Builder {
	builder.query[string(model.Set)] = bson.M{property: value}
	return builder
}

func (builder Builder) UpdateMap(updates map[string]interface{}) Builder {
	set := bson.M{}
	for key, value := range updates {
		set[key] = value
	}

	builder.query[string(model.Set)] = set
	return builder
}

func (builder Builder) NotIn(property string, value []interface{}) Builder {
	builder.query[property] = bson.M{string(model.NotIn): value}
	return builder
}

func (builder Builder) AddToSet(property string, value interface{}) Builder {
	builder.query[string(model.AddToSet)] = bson.M{property: value}
	return builder
}

func (builder Builder) Increment(property string, increment int) Builder {
	builder.query[string(model.Increment)] = bson.M{property: increment}
	return builder
}
