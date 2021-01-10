package di

import (
	"reflect"
	"survey-api/pkg/auth"
	"survey-api/pkg/logger"
	"survey-api/pkg/pagination"
	"survey-api/pkg/poll"
	"survey-api/pkg/user"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	deps          Dependencies
	syncOnceStore []sync.Once
)

type Dependencies struct {
	mongoClient           *mongo.Client
	loggerService         *logger.LoggerService
	tokenService          *auth.TokenService
	cookieService         *auth.CookieService
	authService           *auth.AuthService
	pollPaginationService *poll.PollPaginationService
	authRepo              *auth.AuthRepo
	userRepo              *user.UserRepo
	pollRepo              *poll.PollRepo
	authHandler           *auth.AuthHandler
	pollHandler           *poll.PollHandler
	authMapper            *auth.AuthMapper
	userMapper            *user.UserMapper
	pollMapper            *poll.PollMapper
	paginationMapper      *pagination.PaginationMapper
}

func init() {
	deps = Dependencies{}
	syncOnceStore = make([]sync.Once, reflect.TypeOf(deps).NumField())
}

func Container() *Dependencies {
	return &deps
}

func (deps *Dependencies) dbClient() *mongo.Client {
	if deps.mongoClient == nil {
		syncOnceStore[0].Do(func() {
			deps.mongoClient = createDbClient()
		})
	}

	return deps.mongoClient
}

func (deps *Dependencies) LoggerService() *logger.LoggerService {
	if deps.loggerService == nil {
		syncOnceStore[1].Do(func() {
			loggerService := logger.NewLoggerService()
			deps.loggerService = &loggerService
		})
	}

	return deps.loggerService
}

func (deps *Dependencies) TokenService() *auth.TokenService {
	if deps.tokenService == nil {
		syncOnceStore[2].Do(func() {
			tokenService := auth.NewTokenService()
			deps.tokenService = &tokenService
		})
	}

	return deps.tokenService
}

func (deps *Dependencies) CookieService() *auth.CookieService {
	if deps.cookieService == nil {
		syncOnceStore[3].Do(func() {
			cookieService := auth.NewCookieService()
			deps.cookieService = &cookieService
		})
	}

	return deps.cookieService
}

func (deps *Dependencies) AuthRepo() *auth.AuthRepo {
	if deps.authRepo == nil {
		syncOnceStore[4].Do(func() {
			deps.authRepo = createAuthRepo()
		})
	}

	return deps.authRepo
}

func (deps *Dependencies) UserRepo() *user.UserRepo {
	if deps.userRepo == nil {
		syncOnceStore[5].Do(func() {
			deps.userRepo = createUserRepo()
		})
	}

	return deps.userRepo
}

func (deps *Dependencies) PollRepo() *poll.PollRepo {
	if deps.pollRepo == nil {
		syncOnceStore[6].Do(func() {
			deps.pollRepo = createPollRepo()
		})
	}

	return deps.pollRepo
}

func (deps *Dependencies) AuthHandler() *auth.AuthHandler {
	if deps.authHandler == nil {
		syncOnceStore[7].Do(func() {
			deps.authHandler = createAuthHandler()
		})
	}

	return deps.authHandler
}

func (deps *Dependencies) PaginationMapper() *pagination.PaginationMapper {
	if deps.paginationMapper == nil {
		syncOnceStore[8].Do(func() {
			paginationMapper := pagination.NewPaginationMapper()
			deps.paginationMapper = &paginationMapper
		})
	}

	return deps.paginationMapper
}

func (deps *Dependencies) PollHandler() *poll.PollHandler {
	if deps.pollHandler == nil {
		syncOnceStore[9].Do(func() {
			deps.pollHandler = createPollHandler()
		})
	}

	return deps.pollHandler
}

func (deps *Dependencies) PollPaginationService() *poll.PollPaginationService {
	if deps.pollPaginationService == nil {
		syncOnceStore[10].Do(func() {
			pollPaginationService := poll.NewPollPaginationService()
			deps.pollPaginationService = &pollPaginationService
		})
	}

	return deps.pollPaginationService
}

func (deps *Dependencies) PollMapper() *poll.PollMapper {
	if deps.pollMapper == nil {
		syncOnceStore[11].Do(func() {
			pollMapper := poll.NewPollMapper()
			deps.pollMapper = &pollMapper
		})
	}

	return deps.pollMapper
}

func (deps *Dependencies) AuthMapper() *auth.AuthMapper {
	if deps.authMapper == nil {
		syncOnceStore[12].Do(func() {
			deps.authMapper = createAuthMapper()
		})
	}

	return deps.authMapper
}

func (deps *Dependencies) UserMapper() *user.UserMapper {
	if deps.userMapper == nil {
		syncOnceStore[13].Do(func() {
			userMapper := user.NewUserMapper()
			deps.userMapper = &userMapper
		})
	}

	return deps.userMapper
}

func (deps *Dependencies) AuthService() *auth.AuthService {
	if deps.authService == nil {
		syncOnceStore[14].Do(func() {
			deps.authService = createAuthService()
		})
	}

	return deps.authService
}
