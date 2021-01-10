package auth

import "survey-api/pkg/user"

type AuthMapper struct {
	userMapper *user.UserMapper
}

func NewAuthMapper(userMapper *user.UserMapper) AuthMapper {
	return AuthMapper{userMapper: userMapper}
}

func (mapper AuthMapper) ToUser(userRegister UserRegister) user.User {
	user := user.User{
		FirstName: userRegister.FirstName,
		UserName:  userRegister.UserName,
		Email:     userRegister.Email,
		Password:  userRegister.Password,
		AvatarUrl: userRegister.AvatarUrl,
	}
	user.Init()
	return user
}

func (mapper AuthMapper) ToUserAuth(token string, user user.User) UserAuth {
	return UserAuth{
		Token: token,
		User:  mapper.userMapper.ToUserDetails(user),
	}
}
