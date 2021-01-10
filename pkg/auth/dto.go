package auth

import (
	"errors"
	"survey-api/pkg/user"
	"unicode"

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
	Token string           `json:"token"`
	User  user.UserDetails `json:"user"`
}

func (userRegister UserRegister) Validate() error {
	return validation.ValidateStruct(&userRegister,
		validation.Field(&userRegister.FirstName, validation.Required, validation.Length(2, 20)),
		validation.Field(&userRegister.UserName, validation.Required, validation.Length(3, 20)),
		validation.Field(&userRegister.Email, validation.Required, is.Email),
		validation.Field(&userRegister.Password, validation.Required, validation.By(validatePassword)),
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

func (userLogin UserLogin) Validate() error {
	return validation.ValidateStruct(&userLogin,
		validation.Field(&userLogin.UserName, validation.Required),
		validation.Field(&userLogin.Password, validation.Required),
	)
}
