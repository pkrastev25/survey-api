package refresh

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/user"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(
	logger *logger.Service,
	cookieService *auth.CookieService,
	tokenService *auth.TokenService,
	authRepo *auth.AuthRepo,
	userRepo *user.UserRepo,
	authHandler *auth.AuthHandler,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cookie, err := cookieService.ParseSessionCookie(r)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionId, err := cookieService.ValidateSessionCookie(cookie)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, err := authRepo.FindById(sessionId)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var generatedCookie *http.Cookie
		token := session.Token
		_, err = tokenService.ValidateJwtToken(token)
		if err != nil {
			newCookie, newToken, err := authHandler.RefreshAuth(session)
			if err != nil {
				logger.LogErr(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			generatedCookie = &newCookie
			token = newToken
		}

		user, err := userRepo.FindById(session.UserId.Hex())
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authUser := auth.AuthUser{
			Token: token,
			User:  user.ToClientUser(),
		}
		result, err := json.Marshal(authUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if generatedCookie != nil {
			http.SetCookie(w, generatedCookie)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func init() {
	handler = Init(
		di.Container().Logger,
		di.Container().CookieService,
		di.Container().TokenService,
		di.Container().AuthRepo,
		di.Container().UserRepo,
		di.Container().AuthHandler,
	)
}
