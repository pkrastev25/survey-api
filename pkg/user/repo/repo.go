package repo

import (
	"context"
	"errors"
	"survey-api/pkg/db/query"
	"survey-api/pkg/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (service Service) InsertOne(user model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := service.userCollection().InsertOne(ctx, user)
	if err != nil {
		return user, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return user, errors.New("")
	}

	user.Id = id
	return user, nil
}

func (service Service) FindById(userIdString string) (model.User, error) {
	var user model.User
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return user, err
	}

	return service.FindOne(query.New().Filter("_id", userId))
}

func (service Service) FindOne(query query.Builder) (model.User, error) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := service.userCollection().FindOne(ctx, query.Build())
	defer cancel()
	err := result.Err()
	if err != nil {
		return user, err
	}

	err = result.Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service Service) userCollection() *mongo.Collection {
	return service.client.Database("survey").Collection("user")
}

func (service Service) createUserIndexes() error {
	collection := service.userCollection()
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
