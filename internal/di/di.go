//+build wireinject

package di

import (
	"context"
	"os"
	"survey-api/pkg/auth/handler"
	"survey-api/pkg/logger"
	"survey-api/pkg/mongodb"
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
	dependencies = create()
}

func Container() *Dependencies {
	return dependencies
}

func create() *Dependencies {
	panic(wire.Build(
		wire.Struct(new(logger.Service), "*"),
		createMongodbClient,
		mongodb.New,
		repo.New,
		handler.New,
		packageDependencies,
	))
}

func createMongodbClient() *mongo.Client {
	host := os.Getenv("MONGODB_HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv("MONGODB_PORT")
	if len(port) == 0 {
		port = "27017"
	}

	url := "mongodb://" + host + ":" + port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	clientOptions := options.Client().ApplyURI(url).SetAuth(options.Credential{
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	})
	client, err := mongo.Connect(ctx, clientOptions)
	defer cancel()
	if err != nil {
		panic(err)
	}

	return client
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
