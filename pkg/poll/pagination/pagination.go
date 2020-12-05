package pagination

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"survey-api/pkg/db/pipeline"
	"survey-api/pkg/poll/model"
	paginationmodel "survey-api/pkg/poll/pagination/model"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type Metadata struct {
	Where *paginationmodel.Query
	Limit int
}

func New(query url.Values) (*Metadata, error) {
	whereString := query.Get("where")
	limitString := query.Get("limit")
	metadata := &Metadata{}

	if len(whereString) > 0 {
		where, err := base64.URLEncoding.DecodeString(whereString)
		if err != nil {
			return nil, err
		}

		var query *paginationmodel.Query
		err = json.Unmarshal(where, &query)
		if err != nil {
			return nil, err
		}

		metadata.Where = query
	}

	if len(limitString) > 0 {
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			return nil, err
		}

		if limit <= 0 {
			metadata.Limit = defaultLimit
		} else if limit > maxLimit {
			metadata.Limit = maxLimit
		} else {
			metadata.Limit = limit
		}
	} else {
		metadata.Limit = defaultLimit
	}

	return metadata, nil
}

func (metadata *Metadata) ToDbPagination() bson.M {
	if metadata.Where == nil {
		return bson.M{"next": bson.M{string(pipeline.Limit): metadata.Limit + 1}}
	}

	next, prev := metadata.Where.ToDbQuery()
	next = append(next, bson.M{string(pipeline.Limit): metadata.Limit + 1})
	prev = append(prev, bson.M{string(pipeline.Limit): 1})
	return bson.M{"next": next, "prev": prev}
}

func (metadata *Metadata) CreateLinkHeader(r *http.Request, pagination map[string][]model.Poll) (string, error) {
	linkHeader := ""

	if metadata == nil {
		return linkHeader, nil
	}

	navigation := make(map[string]*Metadata)
	nextResult, exists := pagination["next"]
	if exists {
		if len(nextResult) > metadata.Limit {
			navigation["next"] = metadata.NewFromPoll(&nextResult[len(nextResult)-1])
		}
	}

	prevResult, exists := pagination["prev"]
	if exists {
		if len(prevResult) > 0 {
			navigation["prev"] = metadata.NewFromPoll(&prevResult[0])
		}
	}

	if len(navigation) == 0 {
		return linkHeader, nil
	}

	url := r.Host

	for key, metadata := range navigation {
		where, err := json.Marshal(metadata.Where)
		if err != nil {
			return "", err
		}

		linkHeader += "<" + url + "?limit=" + strconv.Itoa(metadata.Limit) + "?where=" + base64.StdEncoding.EncodeToString(where) + ">; rel=" + key + ", "
	}

	return linkHeader, nil
}

func (metadata *Metadata) NewFromPoll(poll *model.Poll) *Metadata {
	newMetadata := &Metadata{
		Where: metadata.Where.NewFromPoll(poll),
		Limit: metadata.Limit,
	}

	return newMetadata
}

func SetLinkHeader(w http.ResponseWriter, linkHeader string) {
	if len(linkHeader) == 0 {
		return
	}

	w.Header().Add("link", linkHeader)
}
