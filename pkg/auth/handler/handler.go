package handler

import (
	"survey-api/pkg/user/model"
	"survey-api/pkg/user/repo"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo *repo.Service
}

func New(userRepo *repo.Service) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Register(registerUser *model.RegisterUser) (*model.User, error) {
	err := registerUser.Validate()
	if err != nil {
		return nil, err
	}

	user := registerUser.ToUser()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	user, err = s.userRepo.InsertOne(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) VerifyUserCredentials(loginUser *model.LoginUser) (*model.User, error) {
	err := loginUser.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindOne(&model.User{UserName: loginUser.UserName})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
