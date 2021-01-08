package logout

import (
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/db"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(
	logger *logger.Service,
	authHandler *auth.AuthHandler,
	authRepo *auth.AuthRepo,
	tokenService *auth.TokenService,
	cookieService *auth.CookieService,
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

		session, err := authRepo.FindOne(db.NewQueryBuilder().Equal("token", token))
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

		cookie := cookieService.GenerateExpiredCookie()
		http.SetCookie(w, &cookie)
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
