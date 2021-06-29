package polls

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/poll"
)

type deps struct {
	loggerService         *logger.LoggerService
	pollPaginationService *poll.PollPaginationService
	authService           *auth.AuthService
	pollHandler           *poll.PollHandler
	pollMapper            *poll.PollMapper
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

		switch r.Method {
		case http.MethodPost:
			handlePost(w, r, userId, deps)
		case http.MethodGet:
			handleGet(w, r, userId, deps)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, userId string, deps *deps) {
	var pollCreate poll.PollCreate
	err := json.NewDecoder(r.Body).Decode(&pollCreate)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	poll, err := deps.pollHandler.CreatePoll(userId, pollCreate)
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

func handleGet(w http.ResponseWriter, r *http.Request, userId string, deps *deps) {
	query, err := deps.pollPaginationService.ParseQuery(r.URL.Query())
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	polls, paginationNavigation, err := deps.pollHandler.Paginate(query)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pollLists := deps.pollMapper.ToPollLists(polls)
	result, err := json.Marshal(pollLists)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	linkHeader, err := deps.pollPaginationService.CreateLinkHeader(r, paginationNavigation)
	if err != nil {
		deps.loggerService.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deps.pollPaginationService.SetLinkHeader(w, linkHeader)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func init() {
	handler = Init(
		&deps{
			pollPaginationService: di.Container().PollPaginationService(),
			loggerService:         di.Container().LoggerService(),
			authService:           di.Container().AuthService(),
			pollHandler:           di.Container().PollHandler(),
			pollMapper:            di.Container().PollMapper(),
		},
	)
}
