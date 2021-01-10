package auth

import (
	"errors"
	"net/http"
	"survey-api/pkg/user"
)

type AuthService struct {
	tokenService  *TokenService
	cookieService *CookieService
	authRepo      *AuthRepo
}

func NewAuthService(
	tokenService *TokenService,
	cookieService *CookieService,
	authRepo *AuthRepo,
) AuthService {
	return AuthService{
		tokenService:  tokenService,
		cookieService: cookieService,
		authRepo:      authRepo,
	}
}

func (handler AuthService) GenerateAuthUser(user user.User) (http.Cookie, string, error) {
	var cookie http.Cookie
	session, err := handler.authRepo.InsertOne(NewSessionUserId(user.Id))
	if err != nil {
		return cookie, "", err
	}

	return handler.GenerateAuthSession(session)
}

func (handler AuthService) GenerateAuthSession(session Session) (http.Cookie, string, error) {
	var cookie http.Cookie
	token, err := handler.tokenService.GenerateJwtToken(session.UserId.Hex())
	if err != nil {
		return cookie, "", err
	}

	cookie, err = handler.cookieService.GenerateSessionCookie(session)
	return cookie, token, err
}

func (service AuthService) AuthToken(r *http.Request) (string, error) {
	token, err := service.tokenService.ParseJwtToken(r)
	if err != nil {
		return "", errors.New("Malformed token")
	}

	return service.tokenService.ValidateJwtToken(token)
}

func (service AuthService) AuthCookie(r *http.Request) (string, error) {
	cookie, err := service.cookieService.ParseSessionCookie(r)
	if err != nil {
		return "", err
	}

	return service.cookieService.ValidateSessionCookie(cookie)
}
