package register

import (
	"encoding/json"
	"net/http"
	"survey-api/internal/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/user/model"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(logger *logger.Service, authHandler *handler.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var registerUser *model.RegisterUser
		err := json.NewDecoder(r.Body).Decode(&registerUser)
		if err != nil {
			di.Container().Logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := registerUser.ToUser()
		user, err = di.Container().AuthHandler.Register(registerUser)
		if err != nil {
			di.Container().Logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		clientUser := user.ToClientUser()
		result, err := json.Marshal(clientUser)
		if err != nil {
			di.Container().Logger.LogErr(err)
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
