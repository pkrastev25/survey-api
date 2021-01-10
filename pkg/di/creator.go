package di

import (
	"context"
	"errors"
	"os"
	"survey-api/pkg/auth"
	"survey-api/pkg/poll"
	"survey-api/pkg/user"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createDbClient() *mongo.Client {
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
		panic(errors.New("MONGODB_USER is not set"))
	}

	password := os.Getenv("MONGODB_PASSWORD")
	if len(password) == 0 {
		panic(errors.New("MONGODB_PASSWORD is not set"))
	}

	url := "mongodb://" + host + ":" + port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(url).SetAuth(options.Credential{
		Username: user,
		Password: password,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return client
}

func createAuthRepo() *auth.AuthRepo {
	authRepo, err := auth.NewAuthRepo(deps.dbClient())
	if err != nil {
		panic(err)
	}

	return &authRepo
}

func createUserRepo() *user.UserRepo {
	userRepo, err := user.NewUserRepo(deps.dbClient())
	if err != nil {
		panic(err)
	}

	return &userRepo
}

func createPollRepo() *poll.PollRepo {
	pollRepo, err := poll.NewPollRepo(deps.dbClient())
	if err != nil {
		panic(err)
	}

	return &pollRepo
}

func createAuthHandler() *auth.AuthHandler {
	authHandler := auth.NewAuthHandler(
		deps.UserRepo(),
		deps.AuthRepo(),
		deps.TokenService(),
		deps.CookieService(),
		deps.AuthMapper(),
	)
	return &authHandler
}

func createPollHandler() *poll.PollHandler {
	pollhandler := poll.NewPollHandler(
		deps.PollRepo(),
		deps.PaginationMapper(),
		deps.PollMapper(),
	)
	return &pollhandler
}

func createAuthService() *auth.AuthService {
	authService := auth.NewAuthService(
		deps.TokenService(),
		deps.CookieService(),
		deps.AuthRepo(),
	)
	return &authService
}

func createAuthMapper() *auth.AuthMapper {
	authMapper := auth.NewAuthMapper(
		deps.UserMapper(),
	)
	return &authMapper
}
