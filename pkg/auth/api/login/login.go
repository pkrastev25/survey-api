package login

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

func Init(logger *logger.Service, authHandler *auth.AuthHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var loginUser user.LoginUser
		err := json.NewDecoder(r.Body).Decode(&loginUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := authHandler.VerifyUserCredentials(loginUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie, token, err := authHandler.GenerateAuth(user)
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

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func init() {
	handler = Init(
		di.Container().Logger,
		di.Container().AuthHandler,
	)
}
