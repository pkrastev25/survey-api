package user

import (
	"survey-api/pkg/validator"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserDetails struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"avatar_url"`
	Created   string `json:"created"`
}

type UserModify struct {
	FirstName   string `json:"first_name"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	AvatarUrl   string `json:"avatar_url"`
}

func (userModify UserModify) Validate() error {
	return validation.ValidateStruct(&userModify,
		validator.RulesFirstName(&userModify.FirstName),
		validator.RulesPassword(&userModify.OldPassword),
		validator.RulesPassword(&userModify.NewPassword),
		validator.RulesAvatarUrl(&userModify.AvatarUrl),
	)
}
