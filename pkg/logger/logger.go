package logger

import (
	"log"
)

type LoggerService struct{}

func NewLoggerService() LoggerService {
	return LoggerService{}
}

func (service LoggerService) Log(msg string) {
	log.Println(msg)
}

func (service LoggerService) LogErr(err error) {
	log.Println(err)
}
