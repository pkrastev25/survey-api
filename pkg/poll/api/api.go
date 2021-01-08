package api

import (
	"encoding/json"
	"net/http"
	"survey-api/pkg/auth"
	"survey-api/pkg/di"
	"survey-api/pkg/logger"
	"survey-api/pkg/poll"
)

type dependencies struct {
	logger                *logger.Service
	authHandler           *auth.AuthHandler
	pollHandler           *poll.PollHandler
	pollPaginationHandler *poll.PollPaginationHandler
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

		switch r.Method {
		case http.MethodPost:
			handlePost(w, r, userId, deps)
		case http.MethodGet:
			handleGet(w, r, userId, deps)
		case http.MethodDelete:
			handleDelete(w, r, userId, deps)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, userId string, deps *dependencies) {
	var createPoll poll.CreatePoll
	err := json.NewDecoder(r.Body).Decode(&createPoll)
	if err != nil {
		deps.logger.LogErr(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	poll, err := deps.pollHandler.CreatePoll(userId, createPoll)
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

func handleGet(w http.ResponseWriter, r *http.Request, userId string, deps *dependencies) {
	query, err := deps.pollPaginationHandler.ParseQuery(r.URL.Query())
	if err != nil {
		deps.logger.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	polls, paginationNavigation, err := deps.pollHandler.Paginate(query)
	if err != nil {
		deps.logger.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pollClients := make([]poll.PollClient, len(polls))
	for index := range polls {
		pollClients[index] = polls[index].ToPollClient()
	}

	result, err := json.Marshal(pollClients)
	if err != nil {
		deps.logger.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	linkHeader, err := deps.pollPaginationHandler.CreateLinkHeader(r, paginationNavigation)
	if err != nil {
		deps.logger.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deps.pollPaginationHandler.SetLinkHeader(w, linkHeader)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func handleDelete(w http.ResponseWriter, r *http.Request, userId string, deps *dependencies) {
	pollId := r.URL.Query().Get("id")
	if len(pollId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := deps.pollHandler.DeletePoll(userId, pollId)
	if err != nil {
		deps.logger.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func init() {
	handler = Init(
		&dependencies{
			logger:                di.Container().Logger,
			authHandler:           di.Container().AuthHandler,
			pollHandler:           di.Container().PollHandler,
			pollPaginationHandler: di.Container().PollPaginationHandler,
		},
	)
}
