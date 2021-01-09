package auth

import (
	"errors"
	"net/http"
	"survey-api/pkg/db"
	"survey-api/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo      *user.UserRepo
	authRepo      *AuthRepo
	tokenService  *TokenService
	cookieService *CookieService
}

func NewAuthHandler(
	userRepo *user.UserRepo,
	authRepo *AuthRepo,
	tokenService *TokenService,
	cookieService *CookieService,
) AuthHandler {
	return AuthHandler{
		userRepo:      userRepo,
		authRepo:      authRepo,
		tokenService:  tokenService,
		cookieService: cookieService,
	}
}

func (handler AuthHandler) Register(registerUser user.RegisterUser) (user.User, error) {
	var user user.User
	err := registerUser.Validate()
	if err != nil {
		return user, err
	}

	user = registerUser.ToUser()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashedPassword)
	user, err = handler.userRepo.InsertOne(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (handler AuthHandler) VerifyUserCredentials(loginUser user.LoginUser) (user.User, error) {
	var user user.User
	err := loginUser.Validate()
	if err != nil {
		return user, err
	}

	user, err = handler.userRepo.FindOne(db.NewQueryBuilder().Equal("user_name", loginUser.UserName))
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (handler AuthHandler) AuthToken(r *http.Request) (string, error) {
	token, err := handler.tokenService.ParseJwtToken(r)
	if err != nil {
		return "", errors.New("Malformed token")
	}

	userId, err := handler.tokenService.ValidateJwtToken(token)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (handler AuthHandler) GenerateAuth(user user.User) (http.Cookie, string, error) {
	var cookie http.Cookie
	session := Session{
		UserId: user.Id,
	}
	token, err := handler.tokenService.GenerateJwtToken(session.UserId.Hex())
	if err != nil {
		return cookie, "", err
	}

	session.Token = token
	session, err = handler.authRepo.InsertOne(session)
	if err != nil {
		return cookie, "", err
	}

	cookie, err = handler.cookieService.GenerateSessionCookie(session)
	return cookie, token, err
}

func (handler AuthHandler) RefreshAuth(session Session) (http.Cookie, string, error) {
	var cookie http.Cookie
	token, err := handler.tokenService.GenerateJwtToken(session.UserId.Hex())
	if err != nil {
		return cookie, "", err
	}

	filter := db.NewQueryBuilder().Equal("_id", session.Id)
	updates := db.NewQueryBuilder().Set("token", token)
	session, err = handler.authRepo.UpdateOne(filter, updates)
	if err != nil {
		return cookie, "", err
	}

	cookie, err = handler.cookieService.GenerateSessionCookie(session)
	return cookie, token, err
}
