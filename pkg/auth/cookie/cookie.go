package cookie

import (
	"errors"
	"net/http"
	"os"
	"time"

	sessionmodel "survey-api/pkg/auth/model"

	"github.com/gorilla/securecookie"
)

const (
	cookieName          = "survey-session"
	cookieValidityHours = time.Hour * time.Duration(12)
	cookiePath          = "/token/refresh"
)

type Service struct {
}

type cookieStore struct {
	SessionId string `json:"session_id"`
}

func New() *Service {
	return &Service{}
}

func (s *Service) ParseSessionCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(cookieName)
}

func (s *Service) GenerateSessionCookie(session *sessionmodel.Session) (*http.Cookie, error) {
	sessionKey := os.Getenv("SESSION_KEY")
	if len(sessionKey) == 0 {
		return nil, errors.New("SESSION_KEY is not set")
	}

	secureCookie := securecookie.New([]byte(sessionKey), nil)
	cookieStore := &cookieStore{SessionId: session.Id.Hex()}
	encodedValue, err := secureCookie.Encode(cookieName, cookieStore)
	if err != nil {
		return nil, err
	}

	maxAge := int(time.Now().Add(cookieValidityHours).UTC().Unix())
	cookie := s.generateCookie(encodedValue, maxAge)
	return cookie, nil
}

func (s *Service) ValidateSessionCookie(sessionCookie *http.Cookie) (string, error) {
	sessionKey := os.Getenv("SESSION_KEY")
	if len(sessionKey) == 0 {
		return "", errors.New("SESSION_KEY is not set")
	}

	secureCookie := securecookie.New([]byte(sessionKey), nil)
	cookieStore := &cookieStore{}
	err := secureCookie.Decode(cookieName, sessionCookie.Value, &cookieStore)
	if err != nil {
		return "", err
	}

	return cookieStore.SessionId, nil
}

func (s *Service) GenerateExpiredCookie() *http.Cookie {
	return s.generateCookie("", -1)
}

func (s *Service) generateCookie(value string, maxAge int) *http.Cookie {
	cookie := &http.Cookie{
		Name:     cookieName,
		Path:     cookiePath,
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if len(value) != 0 {
		cookie.Value = value
	}

	return cookie
}
