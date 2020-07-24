package handler

import (
	"errors"
	"net/http"
	"survey-api/pkg/auth/cookie"
	authmodel "survey-api/pkg/auth/model"
	authrepo "survey-api/pkg/auth/repo"
	"survey-api/pkg/auth/token"
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

func (s *Service) Register(registerUser *usermodel.RegisterUser) (*usermodel.User, error) {
	err := registerUser.Validate()
	if err != nil {
		return nil, err
	}

	user := registerUser.ToUser()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user, err = s.userRepo.InsertOne(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) VerifyUserCredentials(loginUser *usermodel.LoginUser) (*usermodel.User, error) {
	err := loginUser.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindOne(&usermodel.User{UserName: loginUser.UserName})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AuthToken(r *http.Request) (string, error) {
	token, err := s.tokenService.ParseJwtToken(r)
	if err != nil {
		return "", errors.New("Malformed token")
	}

	userId, err := s.tokenService.ValidateJwtToken(token)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (s *Service) GenerateAuth(user *usermodel.User) (*http.Cookie, string, error) {
	session := &authmodel.Session{
		UserId: user.Id,
	}
	sessionOperation := func(session *authmodel.Session) (*authmodel.Session, error) {
		return s.authRepo.InsertOne(session)
	}

	return s.generateAuthPair(session, sessionOperation)
}

func (s *Service) RefreshAuth(session *authmodel.Session) (*http.Cookie, string, error) {
	sessionOperation := func(session *authmodel.Session) (*authmodel.Session, error) {
		return s.authRepo.ReplaceOne(session)
	}

	return s.generateAuthPair(session, sessionOperation)
}

func (s *Service) generateAuthPair(
	session *authmodel.Session,
	sessionOperation func(*authmodel.Session) (*authmodel.Session, error),
) (*http.Cookie, string, error) {
	token, err := s.tokenService.GenerateJwtToken(session.UserId.Hex())
	if err != nil {
		return nil, "", err
	}

	session.Token = token
	session, err = sessionOperation(session)
	if err != nil {
		return nil, "", err
	}

	cookie, err := s.cookieService.GenerateSessionCookie(session)
	if err != nil {
		return nil, "", err
	}

	return cookie, token, nil
}
