package repo

import (
	"context"
	dbpipeline "survey-api/pkg/db/pipeline"
	"survey-api/pkg/db/query"
	"survey-api/pkg/dtime"
	"survey-api/pkg/poll/model"
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
	err := repo.createPollIndexes()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (service Service) InsertOne(poll model.Poll) (model.Poll, error) {
	poll.LastModified = dtime.DateTimeNow()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := service.pollCollection().InsertOne(ctx, poll)
	if err != nil {
		return poll, err
	}

	poll.Id = result.InsertedID.(primitive.ObjectID)
	return poll, nil
}

func (service Service) FindById(pollIdString string) (model.Poll, error) {
	var poll model.Poll
	pollId, err := primitive.ObjectIDFromHex(pollIdString)
	if err != nil {
		return poll, err
	}

	return service.FindOne(query.New().Filter("_id", pollId))
}

func (service Service) FindOne(query query.Builder) (model.Poll, error) {
	var poll model.Poll
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := service.pollCollection().FindOne(ctx, query.Build())
	err := result.Err()
	if err != nil {
		return poll, err
	}

	err = result.Decode(&poll)
	if err != nil {
		return poll, err
	}

	return poll, nil
}

func (service Service) PaginateQuery(pipeline dbpipeline.Builder) ([]map[string][]model.Poll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cursor, err := service.pollCollection().Aggregate(ctx, pipeline.Build())
	if err != nil {
		return nil, err
	}

	var result []map[string][]model.Poll
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = cursor.All(ctx, &result)
	return result, err
}

func (service Service) UpdateOne(filter query.Builder, updates query.Builder) (model.Poll, error) {
	var poll model.Poll
	updates.Update("last_modified", dtime.DateTimeNow())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := service.pollCollection().FindOneAndUpdate(ctx, filter.Build(), updates.Build(), options)
	err := result.Err()
	if err != nil {
		return poll, err
	}

	err = result.Decode(&poll)
	if err != nil {
		return poll, err
	}

	return poll, nil
}

func (service Service) DeleteOne(filter query.Builder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := service.pollCollection().FindOneAndDelete(ctx, filter.Build())
	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}

func (service Service) pollCollection() *mongo.Collection {
	return service.client.Database("survey").Collection("poll")
}

func (service Service) createPollIndexes() error {
	collection := service.pollCollection()
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
