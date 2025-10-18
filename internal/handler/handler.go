package handler

import (
	"fmt"
	"log"
	"strings"
	"tgbot/internal/fetcher"
)

const token = "8279903094:AAHzqWq_Xx6-CqYLfa9aiedrPx3FJx_sFC4"

type Handler struct {
	fetcher *fetcher.Fetcher
}

func NewHandler(fetcher *fetcher.Fetcher) *Handler {
	return &Handler{
		fetcher: fetcher,
	}
}

func (h *Handler) HandleMessage(chatID int, text string) {
	// Убираем лишние пробелы
	text = strings.TrimSpace(text)

	// Определяем команду
	switch {
	case text == "/start":
		h.handleStart(chatID)
	case strings.HasPrefix(text, "/find"):
		h.handleFind(chatID, text)
	case strings.HasPrefix(text, "/subscribe"):
		// handleSubscribe(chatID, text)
	// case text == "/help":
	// 	handleHelp(chatID)
	default:
		h.handleUnknown(chatID)
	}
}

func (h *Handler) handleStart(chatId int) {
	// сохраняем в бд user и chat id
}

func (h *Handler) handleFind(chatId int, text string) {

	query := strings.TrimPrefix(text, "/find")

	vacancies, err := h.fetcher.Vacancies(query)
	if err != nil {
		log.Printf("Error getting vacancies: %v", err)
		h.fetcher.SendMessage(token, chatId, "Ошибка при поиске вакансий")
	} else {
		for _, vac := range vacancies {
			msg := fmt.Sprintf("%s\n%s\nЗП: %d-%d\n%s",
				vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)
			h.fetcher.SendMessage(token, chatId, msg)
		}
	}

}

func (h *Handler) handleUnknown(chatId int) {
	// отрпавляем что команда не та
}
