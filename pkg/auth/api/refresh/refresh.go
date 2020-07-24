package refresh

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth/cookie"
	authhandler "survey-api/pkg/auth/handler"
	authmodel "survey-api/pkg/auth/model"
	authrepo "survey-api/pkg/auth/repo"
	"survey-api/pkg/auth/token"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	userrepo "survey-api/pkg/user/repo"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(
	logger *logger.Service,
	cookieService *cookie.Service,
	tokenService *token.Service,
	authRepo *authrepo.Service,
	userRepo *userrepo.Service,
	authHandler *authhandler.Service,
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

		var newCookie *http.Cookie
		token := session.Token
		_, err = tokenService.ValidateJwtToken(token)
		if err != nil {
			newCookie, token, err = authHandler.RefreshAuth(session)
			if err != nil {
				logger.LogErr(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		user, err := userRepo.FindById(session.UserId.Hex())
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authUser := &authmodel.AuthUser{
			Token: token,
			User:  user.ToClientUser(),
		}
		result, err := json.Marshal(authUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if newCookie != nil {
			http.SetCookie(w, cookie)
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
