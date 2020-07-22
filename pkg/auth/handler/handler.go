package handler

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"survey-api/pkg/user/model"
	"survey-api/pkg/user/repo"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtTokenValidityMins = time.Minute * time.Duration(10)
	jwtHeader            = "Authorization"
	jwtHeaderValue       = "Bearer "
)

type Service struct {
	userRepo *repo.Service
}

func New(userRepo *repo.Service) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Register(registerUser *model.RegisterUser) (*model.User, error) {
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

func (s *Service) VerifyUserCredentials(loginUser *model.LoginUser) (*model.User, error) {
	err := loginUser.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindOne(&model.User{UserName: loginUser.UserName})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GenerateJwtToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Id.Hex(),
		ExpiresAt: time.Now().Add(jwtTokenValidityMins).Unix(),
	})
	jwtKey := os.Getenv("JWT_KEY")
	if len(jwtKey) == 0 {
		return "", errors.New("JWT_KEY is not set")
	}

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) RequireAuth(r *http.Request) (string, error) {
	parsedToken := strings.Split(r.Header.Get(jwtHeader), jwtHeaderValue)
	if len(parsedToken) != 2 {
		return "", errors.New("Malformed token")
	}

	tokenString := strings.TrimSpace(parsedToken[1])
	return s.ValidateJwtToken(tokenString)
}

func (s *Service) ValidateJwtToken(tokenString string) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")
	if len(jwtKey) == 0 {
		return "", errors.New("JWT_KEY is not set")
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method: " + token.Method.Alg())
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	if len(claims.Subject) == 0 {
		return "", errors.New("Malformed token")
	}

	return claims.Subject, nil
}
