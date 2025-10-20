package sender

import (
	"context"
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

func (s *Sender) Start(ctx context.Context) {
	log.Println("sender is start")

	go s.checkSendMessage(ctx)

}

func (s *Sender) checkSendMessage(ctx context.Context) {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("sender stopped by context")
			return
		case <-ticker.C:
			go s.sendMessages(ctx)
		}
	}
}

func (s *Sender) sendMessages(ctx context.Context) {
	log.Println("sender send message")

	subscriptions, err := s.service.Subscriptions()
	if err != nil {
		log.Println("error fetching subscriptions:", err)
		return
	}

	for _, sub := range subscriptions {
		select {
		case <-ctx.Done():
			log.Println("sendMessages canceled by context")
			return
		default:
		}

		vacancies, err := s.service.SearchVacancies(sub.SearchText, int(sub.TelegramID))
		if err != nil {
			log.Println("error searching vacancies:", err)
			continue
		}

		for _, vac := range vacancies {
			msg := fmt.Sprintf("Подписка: %s\n%s\n%s\nЗП: %d-%d\n%s",
				sub.SearchText, vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)

			s.bot.SendMessage(int(sub.TelegramID), msg)
		}
	}
}
