package userid

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/urlpath"
	"survey-api/pkg/user"
)

const (
	ApiPath = "/users/{id}"
)

type deps struct {
	urlParser     *urlpath.UrlParser
	loggerService *logger.LoggerService
	authService   *auth.AuthService
	userRepo      *user.UserRepo
	userHandler   *user.UserHandler
	userMapper    *user.UserMapper
}

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(
	deps *deps,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := deps.authService.AuthToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params, err := deps.urlParser.ParseParams(r.URL.Path, ApiPath)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGet(w, r, params, userId, deps)
		case http.MethodPut:
			handlePut(w, r, params, userId, deps)
		case http.MethodDelete:
			handleDelete(w, r, params, userId, deps)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, params map[string]string, callerId string, deps *deps) {
	userId := params["id"]
	if len(userId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := deps.userHandler.GetUserById(callerId, userId)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userDetails := deps.userMapper.ToUserDetails(user)
	result, err := json.Marshal(userDetails)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func handlePut(w http.ResponseWriter, r *http.Request, params map[string]string, callerId string, deps *deps) {
	userId := params["id"]
	if len(userId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userModify user.UserModify
	err := json.NewDecoder(r.Body).Decode(&userModify)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := deps.userHandler.ModifyUser(callerId, userId, userModify)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userDetails := deps.userMapper.ToUserDetails(user)
	result, err := json.Marshal(userDetails)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func handleDelete(w http.ResponseWriter, r *http.Request, params map[string]string, callerId string, deps *deps) {
	userId := params["id"]
	if len(userId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := deps.userHandler.DeleteUser(callerId, userId)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func init() {
	handler = Init(
		&deps{
			urlParser:     di.Container().UrlParser(),
			loggerService: di.Container().LoggerService(),
			authService:   di.Container().AuthService(),
			userRepo:      di.Container().UserRepo(),
			userHandler:   di.Container().UserHandler(),
			userMapper:    di.Container().UserMapper(),
		},
	)
}
