package auth

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

var (
	sessionValiditySeconds = (time.Hour * time.Duration(12)).Seconds()
)

type AuthRepo struct {
	client *mongo.Client
}

func NewAuthRepo(client *mongo.Client) (AuthRepo, error) {
	repo := AuthRepo{client: client}
	err := repo.createUserIndexes()
	return repo, err
}

func (repo AuthRepo) InsertOne(session Session) (Session, error) {
	session.UpdateLastModified()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := repo.sessionCollection().InsertOne(ctx, session)
	if err != nil {
		return session, err
	}

	session.Id = result.InsertedID.(primitive.ObjectID)
	return session, nil
}

func (repo AuthRepo) FindById(sessionIdString string) (Session, error) {
	var session Session
	sessionId, err := primitive.ObjectIDFromHex(sessionIdString)
	if err != nil {
		return session, err
	}

	return repo.FindOne(db.NewQueryBuilder().Equal(db.PropertyId, sessionId))
}

func (repo AuthRepo) FindOne(query db.QueryBuilder) (Session, error) {
	var session Session
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := repo.sessionCollection().FindOne(ctx, query.Build())
	defer cancel()
	err := result.Err()
	if err != nil {
		return session, err
	}

	err = result.Decode(&session)
	return session, err
}

func (repo AuthRepo) Execute(pipeline db.PipelineBuilder) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.sessionCollection().Aggregate(ctx, pipeline.Build())
}

func (repo AuthRepo) UpdateOne(filter db.QueryBuilder, updates db.QueryBuilder) (Session, error) {
	var session Session
	updates.Set(db.PropertyLastModified, dtime.DateTimeNow())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := repo.sessionCollection().FindOneAndUpdate(ctx, filter.Build(), updates.Build(), options)
	err := result.Err()
	if err != nil {
		return session, err
	}

	err = result.Decode(&session)
	return session, err
}

func (repo AuthRepo) DeleteOne(session Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := repo.sessionCollection().FindOneAndDelete(ctx, session)
	return result.Err()
}

func (repo AuthRepo) DeleteMany(query db.QueryBuilder) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := repo.sessionCollection().DeleteMany(ctx, query.Build())
	return int(result.DeletedCount), err
}

func (repo AuthRepo) sessionCollection() *mongo.Collection {
	return repo.client.Database(db.DbSurvey).Collection("session")
}

func (repo AuthRepo) createUserIndexes() error {
	collection := repo.sessionCollection()
	indexes := []mongo.IndexModel{
		{
			Keys: bson.M{"last_modified": 1},
			Options: options.Index().SetExpireAfterSeconds(
				int32(sessionValiditySeconds),
			),
		},
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.Indexes().CreateMany(context, indexes)
	defer cancel()
	return err
}
