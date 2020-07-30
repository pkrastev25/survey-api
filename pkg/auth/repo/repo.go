package repo

import (
	"context"
	"survey-api/pkg/auth/model"
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

func (s *Service) InsertOne(session *model.Session) (*model.Session, error) {
	session.Id = primitive.NewObjectID()
	session.LastModified = primitive.NewDateTimeFromTime(time.Now().UTC())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := s.sessionCollection().InsertOne(ctx, session)
	defer cancel()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) FindById(sessionIdString string) (*model.Session, error) {
	sessionId, err := primitive.ObjectIDFromHex(sessionIdString)
	if err != nil {
		return nil, err
	}

	return s.FindOne(&model.Session{Id: sessionId})
}

func (s *Service) FindOne(sessionFilter *model.Session) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := s.sessionCollection().FindOne(ctx, sessionFilter)
	defer cancel()
	err := result.Err()
	if err != nil {
		return nil, err
	}

	var session *model.Session
	err = result.Decode(&session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) ReplaceOne(session *model.Session) (*model.Session, error) {
	session.LastModified = primitive.NewDateTimeFromTime(time.Now().UTC())
	sessionFilter := &model.Session{Id: session.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := s.sessionCollection().ReplaceOne(ctx, sessionFilter, session)
	defer cancel()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) DeleteOne(session *model.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := s.sessionCollection().DeleteOne(ctx, session)
	defer cancel()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) sessionCollection() *mongo.Collection {
	return s.client.Database("survey").Collection("session")
}

func (s *Service) createUserIndexes() error {
	collection := s.sessionCollection()
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
