package auth

import (
	"context"
	"errors"
	"survey-api/pkg/db"
	"survey-api/pkg/dtime"
	"survey-api/pkg/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo      *user.UserRepo
	authRepo      *AuthRepo
	tokenService  *TokenService
	cookieService *CookieService
	authMapper    *AuthMapper
}

func NewAuthHandler(
	userRepo *user.UserRepo,
	authRepo *AuthRepo,
	tokenService *TokenService,
	cookieService *CookieService,
	authMapper *AuthMapper,
) AuthHandler {
	return AuthHandler{
		userRepo:      userRepo,
		authRepo:      authRepo,
		tokenService:  tokenService,
		cookieService: cookieService,
		authMapper:    authMapper,
	}
}

func (handler AuthHandler) Register(userRegister UserRegister) (user.User, error) {
	var user user.User
	err := userRegister.Validate()
	if err != nil {
		return user, err
	}

	user = handler.authMapper.ToUser(userRegister)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashedPassword)
	return handler.userRepo.InsertOne(user)
}

func (handler AuthHandler) Login(userLogin UserLogin) (user.User, error) {
	var userModel user.User
	err := userLogin.Validate()
	if err != nil {
		return userModel, err
	}

	userModel, err = handler.userRepo.FindOne(db.NewQueryBuilder().Equal(user.PropertyUserName, userLogin.UserName))
	if err != nil {
		return userModel, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(userLogin.Password))
	return userModel, err
}

func (handler AuthHandler) RefreshSession(sessionIdString string) (Session, user.User, error) {
	var session Session
	var userModel user.User
	sessionId, err := primitive.ObjectIDFromHex(sessionIdString)
	if err != nil {
		return session, userModel, err
	}

	pipeline := db.NewPipelineBuilder().MatchStage(db.PropertyId, sessionId).SetStage(db.PropertyLastModified, dtime.DateTimeNow()).LookUpStage(user.CollectionUser, propertyUserId, db.PropertyId, "user")
	cursor, err := handler.authRepo.Execute(pipeline)
	if err != nil {
		return session, userModel, err
	}

	var results = []struct {
		Session Session     `bson:",inline"`
		User    []user.User `bson:"user"`
	}{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = cursor.All(ctx, &results)
	if err != nil {
		return session, userModel, err
	}

	if len(results) <= 0 {
		return session, userModel, errors.New("")
	}

	result := results[0]
	if len(result.User) <= 0 {
		return session, userModel, errors.New("")
	}

	return result.Session, result.User[0], nil
}

func (handler AuthHandler) Logout(userIdString string) error {
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return err
	}

	filter := db.NewQueryBuilder().Equal(propertyUserId, userId)
	_, err = handler.authRepo.DeleteMany(filter)
	return err
}
