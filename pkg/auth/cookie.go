package auth

import (
	"errors"
	"net/http"
	"os"
	"survey-api/pkg/dtime"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	cookieName          = "survey-session"
	cookieValidityHours = time.Hour * time.Duration(12)
	cookiePath          = "/token/refresh"
)

type CookieService struct {
}

type cookieStore struct {
	SessionId string `json:"session_id"`
}

func NewCookieService() CookieService {
	return CookieService{}
}

func (service CookieService) ParseSessionCookie(r *http.Request) (http.Cookie, error) {
	var cookie http.Cookie
	sessionCookie, err := r.Cookie(cookieName)
	if err != nil {
		return cookie, err
	}

	if sessionCookie == nil {
		return cookie, errors.New("")
	}

	return *sessionCookie, nil
}

func (service CookieService) GenerateSessionCookie(session Session) (http.Cookie, error) {
	var cookie http.Cookie
	sessionKey := os.Getenv("SESSION_KEY")
	if len(sessionKey) == 0 {
		return cookie, errors.New("SESSION_KEY is not set")
	}

	secureCookie := securecookie.New([]byte(sessionKey), nil)
	cookieStore := &cookieStore{SessionId: session.Id.Hex()}
	encodedValue, err := secureCookie.Encode(cookieName, cookieStore)
	if err != nil {
		return cookie, err
	}

	expires := dtime.TimeNow().Add(cookieValidityHours)
	cookie = service.generateCookie(encodedValue, expires)
	return cookie, nil
}

func (service CookieService) ValidateSessionCookie(sessionCookie http.Cookie) (string, error) {
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

func (service CookieService) GenerateExpiredCookie() http.Cookie {
	return service.generateCookie("", dtime.NilTime)
}

func (service CookieService) generateCookie(value string, expires time.Time) http.Cookie {
	cookie := http.Cookie{
		Name:     cookieName,
		Path:     cookiePath,
		Expires:  expires,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if expires != dtime.NilTime {
		cookie.MaxAge = int(expires.Unix())
	}

	if len(value) > 0 {
		cookie.Value = value
	}

	return cookie
}
