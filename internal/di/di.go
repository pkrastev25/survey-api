//+build wireinject

package di

import (
	"context"
	"flag"
	"os"
	"strings"
	"survey-api/pkg/logger"
	"survey-api/pkg/mongodb"
	"time"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dependencies struct {
	Logger      *logger.Service
	MongoClient *mongodb.Service
}

func Container() *Dependencies {
	return dependencies
}

func MockContainer(c *Dependencies) {
	if !isTestEnv() {
		panic("You cannot mock dependencies during normal application execution! Use mocks only when testing!")
	}

	dependencies = c
}

var dependencies *Dependencies

func init() {
	if isTestEnv() {
		return
	}

	var err error
	dependencies, err = create()
	if err != nil {
		panic(err)
	}
}

func isTestEnv() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.v=") || strings.HasPrefix(arg, "/_test/") || strings.HasPrefix(arg, ".test") {
			return true
		}
	}

	goEnv := os.Getenv("GO_ENV")
	return flag.Lookup("test.v") != nil || goEnv == "test" || goEnv == "testing"
}

func create() (*Dependencies, error) {
	panic(wire.Build(
		createMongodbClient,
		wire.Struct(new(logger.Service), "*"),
		mongodb.New,
		packageDependencies,
	))
}

func createMongodbClient() (*mongo.Client, error) {
	url := os.Getenv("MONGODB_URL")
	if len(url) == 0 {
		url = "mongodb://localhost:27017"
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func packageDependencies(logger *logger.Service, mongoClient *mongodb.Service) *Dependencies {
	return &Dependencies{Logger: logger, MongoClient: mongoClient}
}
