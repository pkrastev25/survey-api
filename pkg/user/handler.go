package user

import (
	"errors"
	"survey-api/pkg/crypt"
	"survey-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	cryptService *crypt.CryptService
	userRepo     *UserRepo
}

func NewUserHandler(
	cryptService *crypt.CryptService,
	userRepo *UserRepo,
) UserHandler {
	return UserHandler{
		cryptService: cryptService,
		userRepo:     userRepo,
	}
}

func (handler UserHandler) GetUserById(callerIdString string, userIdString string) (User, error) {
	var user User
	userId, err := handler.verifyUserAndCallerMatch(callerIdString, userIdString)
	if err != nil {
		return user, err
	}

	filter := db.NewQueryBuilder().Equal(db.PropertyId, userId)
	return handler.userRepo.FindOne(filter)
}

func (handler UserHandler) verifyUserAndCallerMatch(callerIdString string, userIdString string) (primitive.ObjectID, error) {
	var userId primitive.ObjectID
	callerId, err := primitive.ObjectIDFromHex(callerIdString)
	if err != nil {
		return userId, err
	}

	userId, err = primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return userId, err
	}

	if callerId != userId {
		return userId, errors.New("")
	}

	return userId, nil
}

func (handler UserHandler) ModifyUser(callerIdString string, userIdString string, userModify UserModify) (User, error) {
	var user User
	userId, err := handler.verifyUserAndCallerMatch(callerIdString, userIdString)
	if err != nil {
		return user, err
	}

	err = userModify.Validate()
	if err != nil {
		return user, err
	}

	filter := db.NewQueryBuilder().Equal(db.PropertyId, userId)
	updates := db.NewQueryBuilder()
	if len(userModify.FirstName) > 0 {
		updates = updates.Set(PropertyFirstName, userModify.FirstName)
	}

	if len(userModify.OldPassword) > 0 && len(userModify.NewPassword) > 0 {
		if userModify.OldPassword == userModify.NewPassword {
			return user, errors.New("")
		}

		userModel, err := handler.userRepo.FindOne(filter)
		if err != nil {
			return user, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(userModify.OldPassword))
		if err != nil {
			return user, err
		}

		newPasswordHash, err := handler.cryptService.GeneratePasswordHash(userModify.NewPassword)
		if err != nil {
			return user, err
		}

		updates = updates.Set(PropertyPassword, newPasswordHash)
	}

	if len(userModify.AvatarUrl) > 0 {
		updates = updates.Set(PropertyAvatarUrl, userModify.AvatarUrl)
	}

	return handler.userRepo.UpdateOne(filter, updates)
}

func (handler UserHandler) DeleteUser(callerIdString string, userIdString string) error {
	userId, err := handler.verifyUserAndCallerMatch(callerIdString, userIdString)
	if err != nil {
		return err
	}

	filter := db.NewQueryBuilder().Equal(db.PropertyId, userId)
	return handler.userRepo.DeleteOne(filter)
}
