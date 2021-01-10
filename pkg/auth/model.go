package auth

import (
	"survey-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	propertyUserId = "user_id"
)

type Session struct {
	db.BaseModel `bson:",inline"`
	UserId       primitive.ObjectID `bson:"user_id"`
}

func NewSessionUserId(userId primitive.ObjectID) Session {
	return Session{
		BaseModel: db.NewBaseModel(),
		UserId:    userId,
	}
}
