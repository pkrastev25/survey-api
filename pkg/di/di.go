//+build wireinject

package di

import (
	"context"
	"errors"
	"os"
	"survey-api/pkg/auth/handler"
	"survey-api/pkg/logger"
	"survey-api/pkg/user/repo"
	"time"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dependencies struct {
	Logger      *logger.Service
	AuthHandler *handler.Service
}

var dependencies *Dependencies

func init() {
	deps, err := create()
	if err != nil {
		panic(err)
	}

	dependencies = deps
}

func Container() *Dependencies {
	return dependencies
}

func create() (*Dependencies, error) {
	panic(wire.Build(
		wire.Struct(new(logger.Service), "*"),
		createMongodbClient,
		repo.New,
		handler.New,
		packageDependencies,
	))
}

func createMongodbClient() (*mongo.Client, error) {
	host := os.Getenv("MONGODB_HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv("MONGODB_PORT")
	if len(port) == 0 {
		port = "27017"
	}

	user := os.Getenv("MONGODB_USER")
	if len(user) == 0 {
		return nil, errors.New("MONGODB_USER is not set")
	}

	password := os.Getenv("MONGODB_PASSWORD")
	if len(password) == 0 {
		return nil, errors.New("MONGODB_PASSWORD is not set")
	}

	url := "mongodb://" + host + ":" + port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	clientOptions := options.Client().ApplyURI(url).SetAuth(options.Credential{
		Username: user,
		Password: password,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	defer cancel()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func packageDependencies(
	logger *logger.Service,
	authHandler *handler.Service,
) *Dependencies {
	return &Dependencies{
		Logger:      logger,
		AuthHandler: authHandler,
	}
}