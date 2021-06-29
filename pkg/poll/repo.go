package poll

import (
	"context"
	"survey-api/pkg/db"
	"survey-api/pkg/dtime"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PollRepo struct {
	client *mongo.Client
}

func NewPollRepo(client *mongo.Client) (PollRepo, error) {
	repo := PollRepo{client: client}
	err := repo.createPollIndexes()
	return repo, err
}

func (repo PollRepo) InsertOne(poll Poll) (Poll, error) {
	poll.UpdateLastModified()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := repo.pollCollection().InsertOne(ctx, poll)
	if err != nil {
		return poll, err
	}

	poll.Id = result.InsertedID.(primitive.ObjectID)
	return poll, nil
}

func (repo PollRepo) FindById(pollIdString string) (Poll, error) {
	var poll Poll
	pollId, err := primitive.ObjectIDFromHex(pollIdString)
	if err != nil {
		return poll, err
	}

	return repo.FindOne(db.NewQueryBuilder().Equal(db.PropertyId, pollId))
}

func (repo PollRepo) FindOne(filter db.QueryBuilder) (Poll, error) {
	var poll Poll
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := repo.pollCollection().FindOne(ctx, filter.Build())
	err := result.Err()
	if err != nil {
		return poll, err
	}

	err = result.Decode(&poll)
	return poll, err
}

func (repo PollRepo) Execute(pipeline db.PipelineBuilder) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.pollCollection().Aggregate(ctx, pipeline.Build())
}

func (repo PollRepo) UpdateOne(filter db.QueryBuilder, updates db.QueryBuilder) (Poll, error) {
	var poll Poll
	updates.Set(db.PropertyLastModified, dtime.DateTimeNow())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := repo.pollCollection().FindOneAndUpdate(ctx, filter.Build(), updates.Build(), options)
	err := result.Err()
	if err != nil {
		return poll, err
	}

	err = result.Decode(&poll)
	return poll, err
}

func (repo PollRepo) DeleteOne(filter db.QueryBuilder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := repo.pollCollection().FindOneAndDelete(ctx, filter.Build())
	return result.Err()
}

func (repo PollRepo) pollCollection() *mongo.Collection {
	return repo.client.Database(db.DbSurvey).Collection("poll")
}

func (repo PollRepo) createPollIndexes() error {
	collection := repo.pollCollection()
	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"content": "text"},
		},
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateMany(context, indexes)
	return err
}
