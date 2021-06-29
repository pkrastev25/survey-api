package user

import (
	"survey-api/pkg/dtime"
)

type UserMapper struct {
}

func NewUserMapper() UserMapper {
	return UserMapper{}
}

func (mapper UserMapper) ToUserDetails(user User) UserDetails {
	return UserDetails{
		Id:        user.Id.Hex(),
		FirstName: user.FirstName,
		UserName:  user.UserName,
		AvatarUrl: user.AvatarUrl,
		Created:   dtime.DateTimeToISO(user.Created),
	}
}
