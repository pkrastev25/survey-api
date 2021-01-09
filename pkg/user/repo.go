package user

import (
	"context"
	"errors"
	"survey-api/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	client *mongo.Client
}

func NewUserRepo(client *mongo.Client) (UserRepo, error) {
	repo := UserRepo{client: client}
	err := repo.createUserIndexes()
	if err != nil {
		return repo, err
	}

	return repo, nil
}

func (repo UserRepo) InsertOne(user User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := repo.userCollection().InsertOne(ctx, user)
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

func (repo UserRepo) FindById(userIdString string) (User, error) {
	var user User
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return user, err
	}

	return repo.FindOne(db.NewQueryBuilder().Equal("_id", userId))
}

func (repo UserRepo) FindOne(query db.QueryBuilder) (User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := repo.userCollection().FindOne(ctx, query.Build())
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

func (repo UserRepo) userCollection() *mongo.Collection {
	return repo.client.Database("survey").Collection("user")
}

func (repo UserRepo) createUserIndexes() error {
	collection := repo.userCollection()
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
