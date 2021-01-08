package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"survey-api/pkg/dtime"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	jwtHeader            = "Authorization"
	jwtHeaderValue       = "Bearer "
	jwtTokenValidityMins = time.Minute * time.Duration(10)
)

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (service TokenService) ParseJwtToken(r *http.Request) (string, error) {
	parsedToken := strings.Split(r.Header.Get(jwtHeader), jwtHeaderValue)
	if len(parsedToken) != 2 {
		return "", errors.New("Malformed token")
	}

	tokenString := strings.TrimSpace(parsedToken[1])
	if len(tokenString) == 0 {
		return "", errors.New("Missing token")
	}

	return tokenString, nil
}

func (service TokenService) GenerateJwtToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: dtime.TimeNow().Add(jwtTokenValidityMins).Unix(),
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

func (service TokenService) ValidateJwtToken(tokenString string) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")
	if len(jwtKey) == 0 {
		return "", errors.New("JWT_KEY is not set")
	}

	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
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
