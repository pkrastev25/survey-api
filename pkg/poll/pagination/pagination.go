package pagination

import (
	"net/url"
	"strconv"
	"time"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type PollsQuery struct {
	ID    string
	Next  time.Time
	Prev  time.Time
	Limit int
}

func New(query url.Values) (*PollsQuery, error) {
	idString := query.Get("id")
	nextString := query.Get("next")
	prevString := query.Get("prev")
	limitString := query.Get("limit")
	pollsQuery := &PollsQuery{}

	if len(idString) > 0 {
		pollsQuery.ID = idString
	}

	if len(nextString) > 0 {
		next, err := time.Parse(time.Now().UTC().String(), nextString)
		if err != nil {
			return nil, err
		}

		pollsQuery.Next = next
	}

	if len(prevString) > 0 {
		prev, err := time.Parse(time.Now().UTC().String(), prevString)
		if err != nil {
			return nil, err
		}

		pollsQuery.Prev = prev
	}

	if len(limitString) > 0 {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return nil, err
		}

		if limit <= 0 {
			pollsQuery.Limit = defaultLimit
		} else if limit > maxLimit {
			pollsQuery.Limit = maxLimit
		} else {
			pollsQuery.Limit = limit
		}
	} else {
		pollsQuery.Limit = defaultLimit
	}

	return pollsQuery, nil
}
