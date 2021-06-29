package pollid

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/poll"
	"survey-api/pkg/urlpath"
)

const (
	ApiPath = "/polls/{id}"
)

type deps struct {
	urlParser     *urlpath.UrlParser
	loggerService *logger.LoggerService
	authService   *auth.AuthService
	pollHandler   *poll.PollHandler
	pollMapper    *poll.PollMapper
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

func handleGet(w http.ResponseWriter, r *http.Request, params map[string]string, userId string, deps *deps) {
	pollId := params["id"]
	if len(pollId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	poll, err := deps.pollHandler.GetPollById(pollId)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pollDetails := deps.pollMapper.ToPollDetails(poll)
	result, err := json.Marshal(pollDetails)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func handlePut(w http.ResponseWriter, r *http.Request, params map[string]string, userId string, deps *deps) {
	pollId := params["id"]
	if len(pollId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var pollModify poll.PollModify
	err := json.NewDecoder(r.Body).Decode(&pollModify)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	poll, err := deps.pollHandler.ModifyPoll(pollId, userId, pollModify)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pollDetails := deps.pollMapper.ToPollDetails(poll)
	result, err := json.Marshal(pollDetails)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func handleDelete(w http.ResponseWriter, r *http.Request, params map[string]string, userId string, deps *deps) {
	pollId := params["id"]
	if len(pollId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := deps.pollHandler.DeletePoll(userId, pollId)
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
			pollHandler:   di.Container().PollHandler(),
			pollMapper:    di.Container().PollMapper(),
		},
	)
}
