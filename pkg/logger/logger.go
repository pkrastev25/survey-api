package logger

import (
	"log"
)

type Service struct{}

func (s *Service) Log(msg string) {
	log.Println(msg)
}

func (s *Service) LogErr(err error) {
	log.Fatalln(err)
}
