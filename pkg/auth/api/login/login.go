package login

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
)

type deps struct {
	loggerService *logger.LoggerService
	authService   *auth.AuthService
	authHandler   *auth.AuthHandler
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

		var userLogin auth.UserLogin
		err := json.NewDecoder(r.Body).Decode(&userLogin)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := deps.authHandler.Login(userLogin)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie, token, err := deps.authService.GenerateAuthUser(user)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userAuth := deps.authMapper.ToUserAuth(token, user)
		result, err := json.Marshal(userAuth)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func init() {
	handler = Init(
		&deps{
			loggerService: di.Container().LoggerService(),
			authService:   di.Container().AuthService(),
			authHandler:   di.Container().AuthHandler(),
			authMapper:    di.Container().AuthMapper(),
		},
	)
}
