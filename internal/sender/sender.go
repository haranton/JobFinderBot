package sender

import (
	"log"
	"tgbot/internal/service"
	"time"
)

type Sender struct {
	service *service.Service
}

func NewSender(service *service.Service) *Sender {
	return &Sender{
		service: service,
	}
}

func (s *Sender) Start() {
	log.Println("sender is start")
	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for range ticker.C {
			log.Println("sender send message")
			// s.fetcher.Vacancies()
		}
	}()
}
