package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const surveyDbName = "survey"
const userCollectionName = "user"

type Service struct {
	db *mongo.Client
}

func New(client *mongo.Client) *Service {
	createUserIndexes(client)
	return &Service{db: client}
}

func (s *Service) ContextWithTimeout() (context.Context, context.CancelFunc) {
	return contextWithTimeout()
}

func (s *Service) UserCollection() *mongo.Collection {
	return s.db.Database(surveyDbName).Collection(userCollectionName)
}

func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func createUserIndexes(client *mongo.Client) {
	collection := client.Database(surveyDbName).Collection(userCollectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"user_name": "text"},
			Options: options.Index().SetUnique(true),
		}, {
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	context, cancel := contextWithTimeout()
	_, err := collection.Indexes().CreateMany(context, indexes)
	defer cancel()
	if err != nil {
		panic(err)
	}
}
