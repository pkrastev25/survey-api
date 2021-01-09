package pagination

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	queryPaginate = "paginate"
	querySort     = "sort"
	queryLimit    = "limit"
)

type PaginationService struct {
}

func NewPaginationService() PaginationService {
	return PaginationService{}
}

func (service PaginationService) ParseQuery(queries url.Values) (Query, error) {
	paginateString := queries.Get(queryPaginate)
	sortString := queries.Get(querySort)
	limitString := queries.Get(queryLimit)
	query := NewQuery()

	if len(limitString) > 0 {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return query, err
		}

		query.SetLimit(limit)
	}

	if len(paginateString) > 0 {
		paginate, err := service.parsePaginateQuery(paginateString)
		if err != nil {
			return query, err
		}

		query.SetPaginate(paginate)
		if len(sortString) <= 0 {
			query.SetSort(NewSortPaginate(paginate))
		}
	}

	if len(sortString) > 0 {
		sort, err := service.parseSortQuery(sortString)
		if err != nil {
			return query, err
		}

		query.SetSort(sort)
	}

	return query, nil
}

func (service PaginationService) parsePaginateQuery(query string) (Paginate, error) {
	var paginate Paginate
	values := strings.Split(query, ",")
	if len(values) < 3 {
		return paginate, errors.New("")
	}

	if len(values) > 3 {
		return paginate, errors.New("")
	}

	return NewPaginateAll(values[0], values[1], values[2]), nil
}

func (service PaginationService) parseSortQuery(query string) (Sort, error) {
	var sort Sort
	values := strings.Split(query, ",")
	if len(values) < 2 {
		return sort, errors.New("")
	}

	if len(values) > 2 {
		return sort, errors.New("")
	}

	return NewSortAll(values[0], values[1]), nil
}

func (service PaginationService) CreateLinkHeader(r *http.Request, pagination map[string]Query) string {
	var linkHeader string
	if len(pagination) <= 0 {
		return linkHeader
	}

	protocol := "https://"
	if r.Proto == "HTTP/1.1" {
		protocol = "http://"
	}

	var linkHeaderEntries []string
	url := protocol + r.Host + r.URL.Path
	for navigation, query := range pagination {
		queries := service.QueryStrings(query)
		linkHeaderEntries = append(linkHeaderEntries, "<"+url+"?"+strings.Join(queries, "&")+">; rel="+navigation)
	}

	return strings.Join(linkHeaderEntries, ",")
}

func (service PaginationService) QueryStrings(query Query) []string {
	var queries []string
	paginate := service.paginateQueryString(query.Paginate())
	if len(paginate) > 0 {
		queries = append(queries, paginate)
	}

	sort := service.sortQueryString(query.Sort())
	if len(sort) > 0 {
		queries = append(queries, sort)
	}

	limit := service.limitQueryString(query.Limit())
	if len(limit) > 0 {
		queries = append(queries, limit)
	}

	return queries
}

func (service PaginationService) paginateQueryString(paginate Paginate) string {
	values := []string{string(paginate.Property()), string(paginate.Operation()), paginate.Value()}
	return queryPaginate + "=" + strings.Join(values, ",")
}

func (service PaginationService) sortQueryString(sort Sort) string {
	values := []string{string(sort.Property()), string(sort.Order())}
	return querySort + "=" + strings.Join(values, ",")
}

func (service PaginationService) limitQueryString(limit int) string {
	if limit <= 0 {
		return ""
	}

	return queryLimit + "=" + strconv.Itoa(limit)
}

func (service PaginationService) SetLinkHeader(w http.ResponseWriter, linkHeader string) {
	if len(linkHeader) <= 0 {
		return
	}

	w.Header().Add("link", linkHeader)
}
