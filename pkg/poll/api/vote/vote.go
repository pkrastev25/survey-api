package vote

import (
	"encoding/json"
	"net/http"
	authhandler "survey-api/pkg/auth/handler"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	pollhandler "survey-api/pkg/poll/handler"
	"survey-api/pkg/poll/model"
)

type dependencies struct {
	logger      *logger.Service
	authHandler *authhandler.Service
	pollHandler *pollhandler.Service
}

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init(
	deps *dependencies,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := deps.authHandler.AuthToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var pollVote model.PollVote
		err = json.NewDecoder(r.Body).Decode(&pollVote)
		if err != nil {
			deps.logger.LogErr(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		poll, err := deps.pollHandler.AddPollVote(userId, pollVote)
		if err != nil {
			deps.logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(poll.ToPollClient())
		if err != nil {
			deps.logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func init() {
	handler = Init(
		&dependencies{
			logger:      di.Container().Logger,
			authHandler: di.Container().AuthHandler,
			pollHandler: di.Container().PollHandler,
		},
	)
}
