package user

import (
	"survey-api/pkg/db"
)

const (
	PropertyFirstName = "first_name"
	PropertyUserName  = "user_name"
	PropertyPassword  = "password"
	PropertyAvatarUrl = "avatar_url"
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
