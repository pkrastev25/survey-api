package repo

import (
	"context"
	dbpipeline "survey-api/pkg/db/pipeline"
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

func (s *Service) InsertOne(p *model.Poll) (*model.Poll, error) {
	p.LastModified = dtime.DateTimeNow()

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

func (s *Service) PaginateQuery(pipeline dbpipeline.Builder) ([]map[string][]model.Poll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cursor, err := s.pollCollection().Aggregate(ctx, pipeline.Build())
	if err != nil {
		return nil, err
	}

	var result []map[string][]model.Poll
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) AddVote(pollId primitive.ObjectID, userId primitive.ObjectID, pollOptionIndex string) (*model.Poll, error) {
	updates := bson.M{
		"$set":      bson.M{"last_modified": primitive.NewDateTimeFromTime(time.Now().UTC())},
		"$inc":      bson.M{"options." + pollOptionIndex + ".count": 1},
		"$addToSet": bson.M{"voter_ids": userId},
	}

	return s.updateOne(&model.Poll{Id: pollId}, updates)
}

func (s *Service) UpdateOne(poll *model.Poll) (*model.Poll, error) {
	poll.LastModified = primitive.NewDateTimeFromTime(time.Now().UTC())
	pollFilter := &model.Poll{Id: poll.Id}

	return s.updateOne(pollFilter, bson.M{"$set": poll})
}

func (s *Service) updateOne(pollFilter *model.Poll, updates interface{}) (*model.Poll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := s.pollCollection().FindOneAndUpdate(ctx, pollFilter, updates, options)
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
