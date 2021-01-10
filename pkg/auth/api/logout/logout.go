package logout

import (
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
)

type deps struct {
	loggerService *logger.LoggerService
	authHandler   *auth.AuthHandler
	authRepo      *auth.AuthRepo
	tokenService  *auth.TokenService
	cookieService *auth.CookieService
	authService   *auth.AuthService
	authMapper    *auth.AuthMapper
}

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(deps *deps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		userId, err := deps.authService.AuthToken(r)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = deps.authHandler.Logout(userId)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cookie := deps.cookieService.GenerateExpiredCookie()
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}

func init() {
	handler = Init(
		&deps{
			loggerService: di.Container().LoggerService(),
			authHandler:   di.Container().AuthHandler(),
			authRepo:      di.Container().AuthRepo(),
			tokenService:  di.Container().TokenService(),
			cookieService: di.Container().CookieService(),
			authService:   di.Container().AuthService(),
			authMapper:    di.Container().AuthMapper(),
		},
	)
}
