package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	passwordRegex = `^(?=.{8,32}$)(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9]).*`
)

type RegisterUser struct {
	FirstName string `json:"first_name"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
}

type LoginUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type ClientUser struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"avatar_url"`
}

type User struct {
	Id        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"first_name"`
	UserName  string             `bson:"user_name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	AvatarUrl string             `bson:"avatar_url"`
}

func (u *RegisterUser) ToUser() *User {
	return &User{
		FirstName: u.FirstName,
		UserName:  u.UserName,
		Email:     u.Email,
		Password:  u.Password,
		AvatarUrl: u.AvatarUrl,
	}
}

func (u *User) ToClientUser() *ClientUser {
	return &ClientUser{
		Id:        u.Id.String(),
		FirstName: u.FirstName,
		UserName:  u.UserName,
		AvatarUrl: u.AvatarUrl,
	}
}

func (u RegisterUser) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required, validation.Length(2, 20)),
		validation.Field(&u.UserName, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
	)
}

func (u LoginUser) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserName, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
}
