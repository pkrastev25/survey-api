package auth

import (
	"survey-api/pkg/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id"`
	Token        string             `bson:"token"`
	LastModified primitive.DateTime `bson:"last_modified"`
}

type AuthUser struct {
	Token string          `json:"token"`
	User  user.ClientUser `json:"user"`
}
