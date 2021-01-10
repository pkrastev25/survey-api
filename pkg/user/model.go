package user

import (
	"survey-api/pkg/db"
)

const (
	PropertyUserName = "user_name"
)

type User struct {
	db.BaseModel `bson:",inline"`
	FirstName    string `bson:"first_name"`
	UserName     string `bson:"user_name"`
	Email        string `bson:"email"`
	Password     string `bson:"password"`
	AvatarUrl    string `bson:"avatar_url"`
}

func (user *User) Init() {
	user.BaseModel = db.NewBaseModel()
}
