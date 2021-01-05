package repo

import (
	"context"
	"survey-api/pkg/auth/model"
	"survey-api/pkg/db/pipeline"
	"survey-api/pkg/db/query"
	"survey-api/pkg/dtime"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	sessionValiditySeconds = (time.Hour * time.Duration(12)).Seconds()
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

func (service Service) InsertOne(session model.Session) (model.Session, error) {
	session.LastModified = dtime.DateTimeNow()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := service.sessionCollection().InsertOne(ctx, session)
	if err != nil {
		return session, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return session, err
	}

	session.Id = id
	return session, nil
}

func (service Service) FindById(sessionIdString string) (model.Session, error) {
	var session model.Session
	sessionId, err := primitive.ObjectIDFromHex(sessionIdString)
	if err != nil {
		return session, err
	}

	return service.FindOne(query.New().Filter("_id", sessionId))
}

func (service Service) FindOne(query query.Builder) (model.Session, error) {
	var session model.Session
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := service.sessionCollection().FindOne(ctx, query.Build())
	defer cancel()
	err := result.Err()
	if err != nil {
		return session, err
	}

	err = result.Decode(&session)
	return session, err
}

func (service Service) UpdateOne(filter query.Builder, updates query.Builder) (model.Session, error) {
	var session model.Session
	updates.Update("last_modified", dtime.DateTimeNow())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := service.sessionCollection().FindOneAndUpdate(ctx, filter.Build(), updates.Build(), options)
	err := result.Err()
	if err != nil {
		return session, err
	}

	err = result.Decode(&session)
	return session, err
}

func (service Service) DeleteOne(session model.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := service.sessionCollection().FindOneAndDelete(ctx, session)
	return result.Err()
}

func (service Service) Execute(pipeline pipeline.Builder, resultType interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cursor, err := service.sessionCollection().Aggregate(ctx, pipeline.Build())
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = cursor.All(ctx, &resultType)
	return resultType, err
}

func (service Service) sessionCollection() *mongo.Collection {
	return service.client.Database("survey").Collection("session")
}

func (service Service) createUserIndexes() error {
	collection := service.sessionCollection()
	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"last_modified": 1},
			Options: options.Index().SetExpireAfterSeconds(
				int32(sessionValiditySeconds),
			),
		},
		{
			Keys: bson.M{"token": "text"},
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
