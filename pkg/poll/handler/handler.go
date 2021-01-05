package handler

import (
	"survey-api/pkg/db/pipeline"
	"survey-api/pkg/db/query"
	"survey-api/pkg/poll/model"
	paginationmodel "survey-api/pkg/poll/pagination/model"
	"survey-api/pkg/poll/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	pollRepo *repo.Service
}

func New(pollRepo *repo.Service) *Service {
	return &Service{pollRepo: pollRepo}
}

func (service Service) CreatePoll(userId string, createPoll model.CreatePoll) (model.Poll, error) {
	var poll model.Poll
	err := createPoll.Validate()
	if err != nil {
		return poll, err
	}

	poll, err = createPoll.ToPoll(userId)
	if err != nil {
		return poll, err
	}

	return service.pollRepo.InsertOne(poll)
}

func (service Service) PaginatePolls(query paginationmodel.Query) ([]model.Poll, map[string]paginationmodel.Query, error) {
	err := query.Validate()
	if err != nil {
		return nil, nil, err
	}

	paginationPipeline, err := pipeline.New().Pagination(query)
	if err != nil {
		return nil, nil, err
	}

	paginationPipelineResult, err := service.pollRepo.PaginateQuery(paginationPipeline)
	if err != nil {
		return nil, nil, err
	}

	return paginationPipeline.ParsePagination(query, paginationPipelineResult)
}

func (service Service) AddPollVote(userIdString string, pollVote model.PollVote) (model.Poll, error) {
	var poll model.Poll
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

	filter := query.New().Filter("_id", pollId).NotIn("voter_ids", []interface{}{userId})
	updates := query.New().AddToSet("voter_ids", userId).Increment("options."+pollVote.Index+".count", 1)
	return service.pollRepo.UpdateOne(filter, updates)
}

func (service Service) DeletePoll(userIdString string, pollIdString string) error {
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return err
	}

	pollId, err := primitive.ObjectIDFromHex(pollIdString)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{
		"_id":        pollId,
		"creator_id": userId,
	}
	return service.pollRepo.DeleteOne(query.NewMap(filters))
}
