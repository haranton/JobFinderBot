package sender

import (
	"log"
	"tgbot/internal/fetcher"
	"time"
)

type Sender struct {
	fetcher *fetcher.Fetcher
}

func NewSender(fetcher *fetcher.Fetcher) *Sender {
	return &Sender{
		fetcher: fetcher,
	}
}

func (s *Sender) Start() {
	log.Println("sender is start")
	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for range ticker.C {
			// log.Println("sender send message")
			// s.fetcher.Vacancies()
		}
	}()
}
