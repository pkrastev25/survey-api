package user

type UserDetails struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"user_name"`
	AvatarUrl string `json:"avatar_url"`
}
