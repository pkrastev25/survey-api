package handler

import (
	"errors"
	"survey-api/pkg/poll/model"
	"survey-api/pkg/poll/repo"
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
