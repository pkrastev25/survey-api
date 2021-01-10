package db

import (
	"survey-api/pkg/dtime"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Created      primitive.DateTime `bson:"created"`
	LastModified primitive.DateTime `bson:"last_modified"`
}

func NewBaseModel() BaseModel {
	return BaseModel{
		Created:      dtime.DateTimeNow(),
		LastModified: dtime.DateTimeNow(),
	}
}

func (baseModel *BaseModel) UpdateLastModified() {
	baseModel.LastModified = dtime.DateTimeNow()
}
