package poll

import (
	"net/http"
	"net/url"
	"strings"
	"survey-api/pkg/pagination"
)

const (
	querySearch = "search"
)

type PollPaginationService struct {
	base pagination.PaginationService
}

func NewPollPaginationService() PollPaginationService {
	return PollPaginationService{base: pagination.NewPaginationService()}
}

func (service PollPaginationService) ParseQuery(queries url.Values) (QueryPoll, error) {
	queryPoll := NewQueryPoll()
	query, err := service.base.ParseQuery(queries)
	if err != nil {
		return queryPoll, err
	}

	queryPoll.SetBase(query)
	searchString := queries.Get(querySearch)
	if len(searchString) > 0 {
		queryPoll.SetSearch(searchString)
	}

	return queryPoll, nil
}

func (service PollPaginationService) searchQueryString(search string) string {
	if len(search) <= 0 {
		return ""
	}

	return querySearch + "=" + search
}

func (service PollPaginationService) queryStrings(queryPoll QueryPoll) []string {
	queries := service.base.QueryStrings(queryPoll.base)
	search := service.searchQueryString(queryPoll.Search())
	if len(search) > 0 {
		queries = append(queries, search)
	}

	return queries
}

func (service PollPaginationService) CreateLinkHeader(r *http.Request, pagination map[string]QueryPoll) (string, error) {
	var linkHeader string
	if len(pagination) <= 0 {
		return linkHeader, nil
	}

	protocol := "https://"
	if r.Proto == "HTTP/1.1" {
		protocol = "http://"
	}

	var linkHeaderParts []string
	url := protocol + r.Host + r.URL.Path
	for navigation, query := range pagination {
		queries := service.queryStrings(query)
		linkHeaderParts = append(linkHeaderParts, "<"+url+"?"+strings.Join(queries, "&")+">; rel="+navigation)
	}

	return strings.Join(linkHeaderParts, ","), nil
}

func (service PollPaginationService) SetLinkHeader(w http.ResponseWriter, linkHeader string) {
	if len(linkHeader) <= 0 {
		return
	}

	w.Header().Add("link", linkHeader)
}
