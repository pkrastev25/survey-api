package model

import (
	"errors"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	UserName  string             `bson:"user_name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	AvatarUrl string             `bson:"avatar_url"`
}

func (registerUser RegisterUser) ToUser() User {
	return User{
		FirstName: registerUser.FirstName,
		UserName:  registerUser.UserName,
		Email:     registerUser.Email,
		Password:  registerUser.Password,
		AvatarUrl: registerUser.AvatarUrl,
	}
}

func (user User) ToClientUser() ClientUser {
	return ClientUser{
		Id:        user.Id.Hex(),
		FirstName: user.FirstName,
		UserName:  user.UserName,
		AvatarUrl: user.AvatarUrl,
	}
}

func (registerUser RegisterUser) Validate() error {
	return validation.ValidateStruct(&registerUser,
		validation.Field(&registerUser.FirstName, validation.Required, validation.Length(2, 20)),
		validation.Field(&registerUser.UserName, validation.Required, validation.Length(3, 20)),
		validation.Field(&registerUser.Email, validation.Required, is.Email),
		validation.Field(&registerUser.Password, validation.Required, validation.By(validatePassword)),
	)
}

func validatePassword(value interface{}) error {
	password, ok := value.(string)
	if !ok {
		return errors.New("")
	}

	var (
		hasMinLen    = false
		hasUpperCase = false
		hasLowerCase = false
		hasSpecial   = false
		hasNumber    = false
	)
	if len(password) >= 7 {
		hasMinLen = true
	}

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
			continue
		}

		if unicode.IsUpper(char) {
			hasUpperCase = true
			continue
		}

		if unicode.IsNumber(char) {
			hasNumber = true
			continue
		}

		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
			continue
		}
	}

	if !hasMinLen {
		return errors.New("")
	}

	if !hasLowerCase {
		return errors.New("")
	}

	if !hasUpperCase {
		return errors.New("")
	}

	if !hasNumber {
		return errors.New("")
	}

	if !hasSpecial {
		return errors.New("")
	}

	return nil
}

func (loginUser LoginUser) Validate() error {
	return validation.ValidateStruct(&loginUser,
		validation.Field(&loginUser.UserName, validation.Required),
		validation.Field(&loginUser.Password, validation.Required),
	)
}
