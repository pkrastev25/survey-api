package model

import (
	"survey-api/pkg/user/model"
)

type AuthUser struct {
	Token string            `json:"token"`
	User  *model.ClientUser `json:"user"`
}
