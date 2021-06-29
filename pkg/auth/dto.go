package auth

import (
	"survey-api/pkg/validator"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserRegister struct {
	FirstName string `json:"first_name"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
}

type UserLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserAuth struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"avatar_url"`
	Token     string `json:"token"`
}

func (userRegister UserRegister) Validate() error {
	return validation.ValidateStruct(&userRegister,
		validator.RulesFirstName(&userRegister.FirstName, validation.Required),
		validator.RulesUserName(&userRegister.UserName, validation.Required),
		validation.Field(&userRegister.Email, validation.Required, is.Email),
		validator.RulesPassword(&userRegister.Password, validation.Required),
		validator.RulesAvatarUrl(&userRegister.AvatarUrl, validation.Required),
	)
}

func (userLogin UserLogin) Validate() error {
	return validation.ValidateStruct(&userLogin,
		validator.RulesUserName(&userLogin.UserName, validation.Required),
		validator.RulesPassword(&userLogin.Password, validation.Required),
	)
}
