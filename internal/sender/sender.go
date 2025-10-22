package sender

import (
	"context"
	"fmt"
	"log/slog"
	"tgbot/internal/bot"
	"tgbot/internal/service"
	"time"
)

type Sender struct {
	service *service.Service
	bot     *bot.Bot
	slogger *slog.Logger
}

func NewSender(service *service.Service, bot *bot.Bot, slogger *slog.Logger) *Sender {
	return &Sender{
		service: service,
		bot:     bot,
		slogger: slogger,
	}
}

func (s *Sender) Start(ctx context.Context) {
	s.slogger.Info("sender is start")

	go s.checkSendMessage(ctx)

}

func (s *Sender) checkSendMessage(ctx context.Context) {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.slogger.Info("sender stopped by context")
			return
		case <-ticker.C:
			go s.sendMessages(ctx)
		}
	}
}

func (s *Sender) sendMessages(ctx context.Context) {
	s.slogger.Info("sender send message")

	subscriptions, err := s.service.Subscriptions(ctx)
	if err != nil {
		s.slogger.Info("error fetching subscriptions:", "error", err)
		return
	}

	for _, sub := range subscriptions {
		select {
		case <-ctx.Done():
			s.slogger.Info("sendMessages canceled by context")
			return
		default:
		}

		vacancies, err := s.service.SearchVacancies(ctx, sub.SearchText, int(sub.TelegramID))
		if err != nil {
			s.slogger.Info("error searching vacancies:", "error", err)
			continue
		}

		for _, vac := range vacancies {
			msg := fmt.Sprintf("Подписка: %s\n%s\n%s\nЗП: %d-%d\n%s",
				sub.SearchText, vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)

			s.bot.SendMessage(int(sub.TelegramID), msg)
		}
	}
}
