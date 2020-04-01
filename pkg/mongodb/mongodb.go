package mongodb

import (
	"survey-api/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db     *mongo.Client
	logger *logger.Service
}

func New(client *mongo.Client, logger *logger.Service) *Service {
	service := Service{db: client, logger: logger}
	return &service
}
