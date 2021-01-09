package poll

import (
	"context"
	"survey-api/pkg/db"
	"survey-api/pkg/dtime"
	"survey-api/pkg/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PollHandler struct {
	pollRepo         *PollRepo
	paginationMapper *pagination.PaginationMapper
}

func NewPollHandler(pollRepo *PollRepo, paginationMapper *pagination.PaginationMapper) PollHandler {
	return PollHandler{pollRepo: pollRepo, paginationMapper: paginationMapper}
}

func (handler PollHandler) CreatePoll(userId string, createPoll CreatePoll) (Poll, error) {
	var poll Poll
	err := createPoll.Validate()
	if err != nil {
		return poll, err
	}

	poll, err = createPoll.ToPoll(userId)
	if err != nil {
		return poll, err
	}

	return handler.pollRepo.InsertOne(poll)
}

func (handler PollHandler) Paginate(query QueryPoll) ([]Poll, map[string]QueryPoll, error) {
	err := query.Validate()
	if err != nil {
		return nil, nil, err
	}

	pipeline, err := handler.preparePipeline(query)
	if err != nil {
		return nil, nil, err
	}

	cursor, err := handler.pollRepo.FindMany(pipeline)
	if err != nil {
		return nil, nil, err
	}

	var result []map[string][]Poll
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, nil, err
	}

	return handler.parsePagination(query, result)
}

func (handler PollHandler) preparePipeline(query QueryPoll) (db.PipelineBuilder, error) {
	pipeline := db.NewPipelineBuilder()
	search := query.Search()
	if len(search) > 0 {
		pipeline.TextSearchStage(search)
	}

	paginate := query.base.Paginate()
	paginateDb, err := handler.paginationMapper.PaginateToDb(paginate)
	if err != nil {
		return pipeline, err
	}

	reversePaginate := paginate.CloneReverseDirection()
	reversePaginateDb, err := handler.paginationMapper.ReversePaginateToDb(reversePaginate)
	if err != nil {
		return pipeline, err
	}

	sort := query.base.Sort()
	sortDb := handler.paginationMapper.SortToDb(sort)
	reverseSort := sort.CloneReverseOrder()
	reverseSortDb := handler.paginationMapper.SortToDb(reverseSort)

	facet := map[string]db.FacetBuilder{
		paginate.Direction():        db.NewFacetBuilder().Match(paginateDb).Sort(sortDb).Limit(query.base.Limit() + 1),
		reversePaginate.Direction(): db.NewFacetBuilder().Match(reversePaginateDb).Sort(reverseSortDb).Limit(1),
	}
	return pipeline.FacetStage(facet), nil
}

func (handler PollHandler) parsePagination(query QueryPoll, paginationResult []map[string][]Poll) ([]Poll, map[string]QueryPoll, error) {
	var polls []Poll
	paginationNavigation := make(map[string]QueryPoll)
	facetResult := paginationResult[0]
	sort := query.base.Sort()
	paginate := query.base.Paginate()
	facetStageResult := facetResult[paginate.Direction()]
	if len(facetStageResult) == query.base.Limit()+1 {
		polls = facetStageResult[:len(facetStageResult)-1]
		paginationNavigation[paginate.Direction()] = handler.generateQueryPoll(query, facetStageResult[len(facetStageResult)-1], paginate, sort)
	} else {
		polls = facetStageResult
	}

	reversePaginate := paginate.CloneReverseDirection()
	facetStageResult = facetResult[reversePaginate.Direction()]
	if len(facetStageResult) > 0 {
		paginationNavigation[reversePaginate.Direction()] = handler.generateQueryPoll(query, facetStageResult[0], reversePaginate.CloneReverseOperation(), sort.CloneReverseOrder())
	}

	return polls, paginationNavigation, nil
}

func (handler PollHandler) generateQueryPoll(query QueryPoll, poll Poll, paginate pagination.Paginate, sort pagination.Sort) QueryPoll {
	paginate = paginate.CloneValue(dtime.DateTimeToISO(poll.Created))
	return query.ClonePaginate(paginate).CloneSort(sort)
}

func (handler PollHandler) AddPollVote(userIdString string, pollVote PollVote) (Poll, error) {
	var poll Poll
	err := pollVote.Validate()
	if err != nil {
		return poll, err
	}

	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return poll, err
	}

	pollId, err := primitive.ObjectIDFromHex(pollVote.PollId)
	if err != nil {
		return poll, err
	}

	filter := db.NewQueryBuilder().Equal("_id", pollId).NotIn("voter_ids", []interface{}{userId})
	updates := db.NewQueryBuilder().AddToSet("voter_ids", userId).Increment("options."+pollVote.Index+".count", 1)
	return handler.pollRepo.UpdateOne(filter, updates)
}

func (handler PollHandler) DeletePoll(userIdString string, pollIdString string) error {
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return err
	}

	pollId, err := primitive.ObjectIDFromHex(pollIdString)
	if err != nil {
		return err
	}

	filters := db.NewQueryBuilder().Equal("_id", pollId).Equal("creator_id", userId)
	return handler.pollRepo.DeleteOne(filters)
}
