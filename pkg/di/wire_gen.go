// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"survey-api/pkg/auth/cookie"
	"survey-api/pkg/auth/handler"
	repo2 "survey-api/pkg/auth/repo"
	"survey-api/pkg/auth/token"
	"survey-api/pkg/logger"
	handler2 "survey-api/pkg/poll/handler"
	repo3 "survey-api/pkg/poll/repo"
	"survey-api/pkg/user/repo"
	"time"
)

// Injectors from di.go:

func create() (*Dependencies, error) {
	service := &logger.Service{}
	client, err := createMongodbClient()
	if err != nil {
		return nil, err
	}
	repoService, err := repo.New(client)
	if err != nil {
		return nil, err
	}
	service2, err := repo2.New(client)
	if err != nil {
		return nil, err
	}
	tokenService := &token.Service{}
	cookieService := &cookie.Service{}
	handlerService := handler.New(repoService, service2, tokenService, cookieService)
	service3, err := repo3.New(client)
	if err != nil {
		return nil, err
	}
	service4 := handler2.New(service3)
	diDependencies := packageDependencies(service, handlerService, tokenService, cookieService, service2, repoService, service3, service4)
	return diDependencies, nil
}

// di.go:

type Dependencies struct {
	Logger        *logger.Service
	AuthHandler   *handler.Service
	TokenService  *token.Service
	CookieService *cookie.Service
	AuthRepo      *repo2.Service
	UserRepo      *repo.Service
	PollRepo      *repo3.Service
	PollHandler   *handler2.Service
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

func packageDependencies(logger2 *logger.Service,
	authHandler *handler.Service,
	tokenService *token.Service,
	cookieService *cookie.Service,
	authRepo *repo2.Service,
	userRepo *repo.Service,
	pollRepo *repo3.Service,
	pollHandler *handler2.Service,
) *Dependencies {
	return &Dependencies{
		Logger:        logger2,
		AuthHandler:   authHandler,
		TokenService:  tokenService,
		CookieService: cookieService,
		AuthRepo:      authRepo,
		UserRepo:      userRepo,
		PollRepo:      pollRepo,
		PollHandler:   pollHandler,
	}
}
