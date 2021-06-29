package user

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

const (
	CollectionUser = "user"
)

type UserRepo struct {
	client *mongo.Client
}

func NewUserRepo(client *mongo.Client) (UserRepo, error) {
	repo := UserRepo{client: client}
	err := repo.createUserIndexes()
	return repo, err
}

func (repo UserRepo) InsertOneContext(ctx context.Context, user User) (User, error) {
	result, err := repo.userCollection().InsertOne(ctx, user)
	if err != nil {
		return user, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (repo UserRepo) FindById(userIdString string) (User, error) {
	var user User
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return user, err
	}

	return repo.FindOne(db.NewQueryBuilder().Equal(db.PropertyId, userId))
}

func (repo UserRepo) FindOne(filter db.QueryBuilder) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.FindOneContext(ctx, filter)
}

func (repo UserRepo) FindOneContext(ctx context.Context, filter db.QueryBuilder) (User, error) {
	var user User
	result := repo.userCollection().FindOne(ctx, filter.Build())
	err := result.Err()
	if err != nil {
		return user, err
	}

	err = result.Decode(&user)
	return user, err
}

func (repo UserRepo) UpdateOne(filter db.QueryBuilder, updates db.QueryBuilder) (User, error) {
	var user User
	updates.Set(db.PropertyLastModified, dtime.DateTimeNow())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := repo.userCollection().FindOneAndUpdate(ctx, filter.Build(), updates.Build(), options)
	err := result.Err()
	if err != nil {
		return user, err
	}

	err = result.Decode(&user)
	return user, err
}

func (repo UserRepo) DeleteOne(filter db.QueryBuilder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := repo.userCollection().FindOneAndDelete(ctx, filter.Build())
	return result.Err()
}

func (repo UserRepo) userCollection() *mongo.Collection {
	return repo.client.Database(db.DbSurvey).Collection(CollectionUser)
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
	return err
}
