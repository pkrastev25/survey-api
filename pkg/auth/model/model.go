package model

import (
	"survey-api/pkg/user/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthUser struct {
	Token string            `json:"token"`
	User  *model.ClientUser `json:"user"`
}

type Session struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id,omitempty"`
	Token        string             `bson:"token,omitempty"`
	LastModified primitive.DateTime `bson:"last_modified,omitempty"`
}
