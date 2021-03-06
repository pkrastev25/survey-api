package repo

import (
	"context"
	"survey-api/pkg/poll/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	client *mongo.Client
}

func New(client *mongo.Client) (*Service, error) {
	repo := &Service{client: client}
	err := repo.createPollIndexes()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *Service) InsertOne(p *model.Poll) (*model.Poll, error) {
	p.LastModified = primitive.NewDateTimeFromTime(time.Now().UTC())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := s.pollCollection().InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) FindById(pollIdString string) (*model.Poll, error) {
	pollId, err := primitive.ObjectIDFromHex(pollIdString)
	if err != nil {
		return nil, err
	}

	return s.FindOne(&model.Poll{Id: pollId})
}

func (s *Service) FindOne(pollFilter *model.Poll) (*model.Poll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := s.pollCollection().FindOne(ctx, pollFilter)
	err := result.Err()
	if err != nil {
		return nil, err
	}

	var poll *model.Poll
	err = result.Decode(&poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *Service) UpdateOne(poll *model.Poll) (*model.Poll, error) {
	poll.LastModified = primitive.NewDateTimeFromTime(time.Now())
	pollFilter := &model.Poll{Id: poll.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := s.pollCollection().FindOneAndUpdate(ctx, pollFilter, bson.M{"$set": poll})
	err := result.Err()
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *Service) DeleteOne(poll *model.Poll) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := s.pollCollection().FindOneAndDelete(ctx, poll)
	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) pollCollection() *mongo.Collection {
	return s.client.Database("survey").Collection("poll")
}

func (s *Service) createPollIndexes() error {
	collection := s.pollCollection()
	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"content": "text"},
		},
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateMany(context, indexes)
	if err != nil {
		return err
	}

	return nil
}
