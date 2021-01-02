package pagination

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	paginationmodel "survey-api/pkg/poll/pagination/model"
)

func ParseQuery(urlQuery url.Values) (paginationmodel.Query, error) {
	searchString := urlQuery.Get("search")
	paginateString := urlQuery.Get("paginate")
	sortString := urlQuery.Get("sort")
	limitString := urlQuery.Get("limit")
	query := paginationmodel.NewQuery()

	if len(limitString) > 0 {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return query, err
		}

		query.SetLimit(limit)
	}

	if len(searchString) > 0 {
		query.SetSearch(searchString)
	}

	if len(paginateString) > 0 {
		paginate, err := parsePaginateQuery(paginateString)
		if err != nil {
			return query, err
		}

		query.SetPaginate(paginate)
		if len(sortString) <= 0 {
			sort, err := paginationmodel.NewSortPaginate(paginate)
			if err != nil {
				return query, err
			}

			query.SetSort(sort)
		}
	}

	if len(sortString) > 0 {
		sort, err := parseSortQuery(sortString)
		if err != nil {
			return query, err
		}

		query.SetSort(sort)
	}

	return query, nil
}

func parsePaginateQuery(query string) (paginationmodel.Paginate, error) {
	var paginate paginationmodel.Paginate
	values := strings.Split(query, ",")
	if len(values) < 3 {
		return paginate, errors.New("")
	}

	if len(values) > 3 {
		return paginate, errors.New("")
	}

	paginate, err := paginationmodel.NewPaginateAll(values[0], values[1], values[2])
	if err != nil {
		return paginate, err
	}

	return paginate, nil
}

func parseSortQuery(query string) (paginationmodel.Sort, error) {
	var sort paginationmodel.Sort
	values := strings.Split(query, ",")
	if len(values) < 2 {
		return sort, errors.New("")
	}

	if len(values) > 2 {
		return sort, errors.New("")
	}

	sort, err := paginationmodel.NewSortAll(values[0], values[1])
	if err != nil {
		return sort, err
	}

	return sort, nil
}

func paginateQueryString(paginate paginationmodel.Paginate) string {
	values := []string{string(paginate.Property()), string(paginate.Operation()), paginate.Value()}
	return "paginate=" + strings.Join(values, ",")
}

func sortQueryString(sort paginationmodel.Sort) string {
	values := []string{string(sort.Property()), string(sort.Order())}
	return "sort=" + strings.Join(values, ",")
}

func limitQueryString(limit int) string {
	if limit <= 0 {
		return ""
	}

	return "limit=" + strconv.Itoa(limit)
}

func searchQueryString(search string) string {
	if len(search) <= 0 {
		return ""
	}

	return "search=" + search
}

func CreateLinkHeader(r *http.Request, pagination map[string]paginationmodel.Query) (string, error) {
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
	for navigation, metadata := range pagination {
		var queries []string
		paginate := paginateQueryString(metadata.Paginate())
		if len(paginate) > 0 {
			queries = append(queries, paginate)
		}

		search := searchQueryString(metadata.Search())
		if len(search) > 0 {
			queries = append(queries, search)
		}

		sort := sortQueryString(metadata.Sort())
		if len(sort) > 0 {
			queries = append(queries, sort)
		}

		limit := limitQueryString(metadata.Limit())
		if len(limit) > 0 {
			queries = append(queries, limit)
		}

		linkHeaderParts = append(linkHeaderParts, "<"+url+"?"+strings.Join(queries, "&")+">; rel="+navigation)
	}

	return strings.Join(linkHeaderParts, ","), nil
}

func SetLinkHeader(w http.ResponseWriter, linkHeader string) {
	if len(linkHeader) <= 0 {
		return
	}

	w.Header().Add("link", linkHeader)
}
