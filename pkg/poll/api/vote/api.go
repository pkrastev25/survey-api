package vote

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
	ApiPath = "/polls/{id}/vote"
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

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		params, err := deps.urlParser.ParseParams(r.URL.Path, ApiPath)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pollId := params["id"]
		if len(pollId) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var pollVote poll.PollVote
		err = json.NewDecoder(r.Body).Decode(&pollVote)
		if err != nil {
			deps.loggerService.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		poll, err := deps.pollHandler.AddPollVote(pollId, userId, pollVote)
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
