package auth

import (
	"context"
	"survey-api/pkg/crypt"
	"survey-api/pkg/db"
	"survey-api/pkg/user"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo      *user.UserRepo
	authRepo      *AuthRepo
	cryptService  *crypt.CryptService
	tokenService  *TokenService
	cookieService *CookieService
	authMapper    *AuthMapper
}

func NewAuthHandler(
	userRepo *user.UserRepo,
	authRepo *AuthRepo,
	cryptService *crypt.CryptService,
	tokenService *TokenService,
	cookieService *CookieService,
	authMapper *AuthMapper,
) AuthHandler {
	return AuthHandler{
		userRepo:      userRepo,
		authRepo:      authRepo,
		cryptService:  cryptService,
		tokenService:  tokenService,
		cookieService: cookieService,
		authMapper:    authMapper,
	}
}

func (handler AuthHandler) Register(userRegister UserRegister) (user.User, Session, error) {
	var userModel user.User
	var sessionModel Session
	err := userRegister.Validate()
	if err != nil {
		return userModel, sessionModel, err
	}

	userModel = handler.authMapper.ToUser(userRegister)
	hashedPassword, err := handler.cryptService.GeneratePasswordHash(userModel.Password)
	if err != nil {
		return userModel, sessionModel, err
	}

	userModel.Password = hashedPassword
	userKey := "user"
	sessionKey := "session"
	createUserAndSession := func(context context.Context) (interface{}, error) {
		user, err := handler.userRepo.InsertOneContext(context, userModel)
		if err != nil {
			return nil, err
		}

		session, err := handler.authRepo.InsertOneContext(context, NewSessionUserId(user.Id))
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			userKey:    user,
			sessionKey: session,
		}
		return result, nil
	}
	result, err := handler.authRepo.Transaction(createUserAndSession)
	if err != nil {
		return userModel, sessionModel, err
	}

	resultMap := result.(map[string]interface{})
	userModel = resultMap[userKey].(user.User)
	sessionModel = resultMap[sessionKey].(Session)
	return userModel, sessionModel, nil
}

func (handler AuthHandler) Login(userLogin UserLogin) (user.User, Session, error) {
	var userModel user.User
	var sessionModel Session
	err := userLogin.Validate()
	if err != nil {
		return userModel, sessionModel, err
	}

	filter := db.NewQueryBuilder().Equal(user.PropertyUserName, userLogin.UserName)
	userModel, err = handler.userRepo.FindOne(filter)
	if err != nil {
		return userModel, sessionModel, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(userLogin.Password))
	if err != nil {
		return userModel, sessionModel, err
	}

	sessionModel, err = handler.authRepo.InsertOne(NewSessionUserId(userModel.Id))
	if err != nil {
		return userModel, sessionModel, err
	}

	return userModel, sessionModel, nil
}

func (handler AuthHandler) RefreshSession(sessionIdString string) (user.User, Session, error) {
	var userModel user.User
	var sessionModel Session
	sessionId, err := primitive.ObjectIDFromHex(sessionIdString)
	if err != nil {
		return userModel, sessionModel, err
	}

	userKey := "user"
	sessionKey := "session"
	refreshSession := func(context context.Context) (interface{}, error) {
		sessionFilter := db.NewQueryBuilder().Equal(db.PropertyId, sessionId)
		session, err := handler.authRepo.UpdateOneContext(context, sessionFilter, db.NewQueryBuilder())
		if err != nil {
			return nil, err
		}

		userFilter := db.NewQueryBuilder().Equal(db.PropertyId, session.UserId)
		user, err := handler.userRepo.FindOneContext(context, userFilter)
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			userKey:    user,
			sessionKey: session,
		}
		return result, nil
	}
	result, err := handler.authRepo.Transaction(refreshSession)
	if err != nil {
		return userModel, sessionModel, err
	}

	resultMap := result.(map[string]interface{})
	userModel = resultMap[userKey].(user.User)
	sessionModel = resultMap[sessionKey].(Session)
	return userModel, sessionModel, nil
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
