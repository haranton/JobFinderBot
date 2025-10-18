package handler

import (
	"fmt"
	"log"
	"strings"
	"tgbot/internal/bot"
	"tgbot/internal/service"
)

const token = "8279903094:AAHzqWq_Xx6-CqYLfa9aiedrPx3FJx_sFC4"

type Handler struct {
	bot     *bot.Bot
	service *service.Service
}

func NewHandler(service *service.Service, bot *bot.Bot) *Handler {
	return &Handler{
		service: service,
		bot:     bot,
	}
}

func (h *Handler) HandleMessage(chatID int, text string) {

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

	vacancies, err := h.service.SearchVacancies(query, chatId)
	if err != nil {
		log.Printf("Error getting vacancies: %v", err)
		h.bot.SendMessage(chatId, "Ошибка при отправке сообщения")
	} else {
		for _, vac := range vacancies {
			msg := fmt.Sprintf("%s\n%s\nЗП: %d-%d\n%s",
				vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)
			h.bot.SendMessage(chatId, msg)
		}
	}

}

func (h *Handler) handleUnknown(chatId int) {
	log.Printf("unknown command")
	h.bot.SendMessage(chatId, "Ошибка при отправке сообщения")
}
