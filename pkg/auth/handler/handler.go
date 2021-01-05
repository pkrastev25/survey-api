package handler

import (
	"errors"
	"net/http"
	"survey-api/pkg/auth/cookie"
	authmodel "survey-api/pkg/auth/model"
	authrepo "survey-api/pkg/auth/repo"
	"survey-api/pkg/auth/token"
	"survey-api/pkg/db/query"
	usermodel "survey-api/pkg/user/model"
	userrepo "survey-api/pkg/user/repo"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo      *userrepo.Service
	authRepo      *authrepo.Service
	tokenService  *token.Service
	cookieService *cookie.Service
}

func New(
	userRepo *userrepo.Service,
	authRepo *authrepo.Service,
	tokenService *token.Service,
	cookieService *cookie.Service,
) *Service {
	return &Service{
		userRepo:      userRepo,
		authRepo:      authRepo,
		tokenService:  tokenService,
		cookieService: cookieService,
	}
}

func (service Service) Register(registerUser usermodel.RegisterUser) (usermodel.User, error) {
	var user usermodel.User
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
	user, err = service.userRepo.InsertOne(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service Service) VerifyUserCredentials(loginUser usermodel.LoginUser) (usermodel.User, error) {
	var user usermodel.User
	err := loginUser.Validate()
	if err != nil {
		return user, err
	}

	user, err = service.userRepo.FindOne(query.New().Filter("user_name", loginUser.UserName))
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service Service) AuthToken(r *http.Request) (string, error) {
	token, err := service.tokenService.ParseJwtToken(r)
	if err != nil {
		return "", errors.New("Malformed token")
	}

	userId, err := service.tokenService.ValidateJwtToken(token)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (service Service) GenerateAuth(user usermodel.User) (http.Cookie, string, error) {
	var cookie http.Cookie
	session := authmodel.Session{
		UserId: user.Id,
	}
	token, err := service.tokenService.GenerateJwtToken(session.UserId.Hex())
	if err != nil {
		return cookie, "", err
	}

	session.Token = token
	session, err = service.authRepo.InsertOne(session)
	if err != nil {
		return cookie, "", err
	}

	cookie, err = service.cookieService.GenerateSessionCookie(session)
	return cookie, token, err
}

func (service Service) RefreshAuth(session authmodel.Session) (http.Cookie, string, error) {
	var cookie http.Cookie
	token, err := service.tokenService.GenerateJwtToken(session.UserId.Hex())
	if err != nil {
		return cookie, "", err
	}

	filter := query.New().Filter("_id", session.Id)
	updates := query.New().Update("token", token)
	session, err = service.authRepo.UpdateOne(filter, updates)
	if err != nil {
		return cookie, "", err
	}

	cookie, err = service.cookieService.GenerateSessionCookie(session)
	return cookie, token, err
}
