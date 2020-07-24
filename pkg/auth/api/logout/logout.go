package logout

import (
	"net/http"
	"survey-api/pkg/auth/cookie"
	authhandler "survey-api/pkg/auth/handler"
	authmodel "survey-api/pkg/auth/model"
	authrepo "survey-api/pkg/auth/repo"
	"survey-api/pkg/auth/token"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(
	logger *logger.Service,
	authHandler *authhandler.Service,
	authRepo *authrepo.Service,
	tokenService *token.Service,
	cookieService *cookie.Service,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		token, err := tokenService.ParseJwtToken(r)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = tokenService.ValidateJwtToken(token)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session, err := authRepo.FindOne(&authmodel.Session{Token: token})
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = authRepo.DeleteOne(session)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, cookieService.GenerateExpiredCookie())
		w.WriteHeader(http.StatusOK)
	}
}

func init() {
	handler = Init(
		di.Container().Logger,
		di.Container().AuthHandler,
		di.Container().AuthRepo,
		di.Container().TokenService,
		di.Container().CookieService,
	)
}
