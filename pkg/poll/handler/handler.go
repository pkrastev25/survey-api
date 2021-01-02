package handler

import (
	"errors"
	"strconv"
	"survey-api/pkg/db/pipeline"
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

func (s *Service) CreatePoll(userId string, createPoll *model.CreatePoll) (*model.Poll, error) {
	err := createPoll.Validate()
	if err != nil {
		return nil, err
	}

	poll, err := createPoll.ToPoll(userId)
	if err != nil {
		return nil, err
	}

	poll, err = s.pollRepo.InsertOne(poll)
	if err != nil {
		return nil, err
	}

	return poll, err
}

func (s *Service) PaginatePolls(query paginationmodel.Query) ([]model.Poll, map[string]paginationmodel.Query, error) {
	err := query.Validate()
	if err != nil {
		return nil, nil, err
	}

	paginationPipeline, err := pipeline.New().Pagination(query)
	if err != nil {
		return nil, nil, err
	}

	paginationPipelineResult, err := s.pollRepo.PaginateQuery(paginationPipeline)
	if err != nil {
		return nil, nil, err
	}

	return paginationPipeline.ParsePagination(query, paginationPipelineResult)
}

func (s *Service) AddPollVote(userIdString string, pollVote *model.PollVote) (*model.Poll, error) {
	err := pollVote.Validate()
	if err != nil {
		return nil, err
	}

	poll, err := s.pollRepo.FindById(pollVote.PollId)
	if err != nil {
		return nil, err
	}

	if poll.Visibility != model.Public {
		return nil, errors.New("Cannot vote for this poll")
	}

	index, err := strconv.Atoi(pollVote.Index)
	if err != nil {
		return nil, err
	}

	if index < 0 || index >= len(poll.Options) {
		return nil, errors.New("Index is out of range")
	}

	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return nil, err
	}

	for i := range poll.VoterIds {
		if poll.VoterIds[i] == userId {
			return nil, errors.New("User already voted for this poll")
		}
	}

	poll, err = s.pollRepo.AddVote(poll.Id, userId, pollVote.Index)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *Service) DeletePoll(userId string, pollId string) error {
	poll, err := s.pollRepo.FindById(pollId)
	if err != nil {
		return err
	}

	if poll.OwnerId.Hex() != userId {
		return errors.New("User cannot delete this poll")
	}

	err = s.pollRepo.DeleteOne(poll)
	if err != nil {
		return err
	}

	return nil
}
