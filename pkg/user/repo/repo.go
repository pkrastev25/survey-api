package repo

import (
	"context"
	"survey-api/pkg/mongodb"
	"survey-api/pkg/user/model"
	"time"
)

type Service struct {
	mongodb *mongodb.Service
}

func New(mongodb *mongodb.Service) *Service {
	return &Service{mongodb: mongodb}
}

func (s *Service) InsertOne(u *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := s.mongodb.UserCollection().InsertOne(ctx, u)
	defer cancel()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) FindOne(userFilter *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := s.mongodb.UserCollection().FindOne(ctx, userFilter)
	defer cancel()
	err := result.Err()
	if err != nil {
		return nil, err
	}

	var user *model.User
	err = result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
