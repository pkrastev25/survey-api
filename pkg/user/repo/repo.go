package repo

import (
	"context"
	"survey-api/pkg/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	client *mongo.Client
}

func New(client *mongo.Client) (*Service, error) {
	repo := &Service{client: client}
	err := repo.createUserIndexes()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *Service) InsertOne(u *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := s.userCollection().InsertOne(ctx, u)
	defer cancel()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) FindOne(userFilter *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := s.userCollection().FindOne(ctx, userFilter)
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

func (s *Service) userCollection() *mongo.Collection {
	return s.client.Database("survey").Collection("user")
}

func (s *Service) createUserIndexes() error {
	collection := s.userCollection()
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"user_name": "text"},
			Options: options.Index().SetUnique(true),
		}, {
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.Indexes().CreateMany(context, indexes)
	defer cancel()
	if err != nil {
		return err
	}

	return nil
}
