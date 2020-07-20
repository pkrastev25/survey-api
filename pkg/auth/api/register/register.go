package register

import (
	"encoding/json"
	"net/http"
	authHandler "survey-api/pkg/auth/handler"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	userModel "survey-api/pkg/user/model"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(logger *logger.Service, authHandler *authHandler.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var registerUser *userModel.RegisterUser
		err := json.NewDecoder(r.Body).Decode(&registerUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := registerUser.ToUser()
		user, err = authHandler.Register(registerUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		clientUser := user.ToClientUser()
		result, err := json.Marshal(clientUser)
		if err != nil {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(result)
		w.WriteHeader(http.StatusOK)
	}
}

func init() {
	handler = Init(
		di.Container().Logger,
		di.Container().AuthHandler,
	)
}
