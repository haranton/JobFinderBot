package sender

import (
	"fmt"
	"log"
	"tgbot/internal/bot"
	"tgbot/internal/service"
	"time"
)

type Sender struct {
	service *service.Service
	bot     *bot.Bot
}

func NewSender(service *service.Service, bot *bot.Bot) *Sender {
	return &Sender{
		service: service,
		bot:     bot,
	}
}

func (s *Sender) Start() {
	log.Println("sender is start")
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			log.Println("sender send message")
			subscriptions, err := s.service.Subscriptions()
			if err != nil {
				log.Println(err)
			}

			for _, sub := range subscriptions {
				vakacies, err := s.service.SearchVacancies(sub.SearchText, int(sub.TelegramID))
				if err != nil {
					log.Println(vakacies)
					continue
				}
				for _, vac := range vakacies {
					msg := fmt.Sprintf("%s\n%s\nЗП: %d-%d\n%s",
						vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)
					s.bot.SendMessage(int(sub.TelegramID), msg)
				}
			}
		}
	}()
}
