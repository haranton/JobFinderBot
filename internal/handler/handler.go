package handler

import (
	"fmt"
	"log"
	"strings"
	"tgbot/internal/bot"
	"tgbot/internal/models"
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

func (h *Handler) HandleMessage(userId int, text string) {

	text = strings.TrimSpace(text)

	// Определяем команду
	switch {
	case text == "/start":
		h.handleStart(userId)
	case strings.HasPrefix(text, "/find"):
		h.handleFind(userId, text)
	case strings.HasPrefix(text, "/subscribe"):
		h.handleSubscribe(userId, text)
	// case text == "/help":
	// 	handleHelp(userId)
	default:
		h.handleUnknown(userId)
	}
}

func (h *Handler) handleSubscribe(userId int, subscribeCommandText string) {

	subscribeText := strings.TrimPrefix(subscribeCommandText, "/subscribe")
	// сохраняем в бд user и chat id
	subscribe, err := h.service.RegisterSubscribe(userId, subscribeText)

	if err != nil {
		log.Println(err)
		h.bot.SendMessage(userId, "Ошибка регистрации подписки")
		return
	}

	msg := fmt.Sprintf("Подписка успешно зарегистрирована, данные подписки: %v", subscribe)
	if err = h.bot.SendMessage(userId, msg); err != nil {
		log.Println("error send message")
	}

}

func (h *Handler) handleStart(userId int) {
	// сохраняем в бд user и chat id
	user, err := h.service.RegisterUser(userId)
	if err != nil && user != (models.User{}) {
		log.Println(err)
		h.bot.SendMessage(userId, "Пользователь уже зарегистрирован")
		return
	}
	if err != nil {
		log.Println(err)
		h.bot.SendMessage(userId, "Ошибка регистрации пользователя")
		return
	}

	msg := fmt.Sprintf("пользователь успешно зарегистрирован, данные пользователя: %v", user)
	if err = h.bot.SendMessage(userId, msg); err != nil {
		log.Println("error send message")
	}

}

func (h *Handler) handleFind(userId int, text string) {

	query := strings.TrimPrefix(text, "/find")

	vacancies, err := h.service.SearchVacancies(query, userId)
	if err != nil {
		log.Printf("Error getting vacancies: %v", err)
		h.bot.SendMessage(userId, "Ошибка при отправке сообщения")
	} else {
		for _, vac := range vacancies {
			msg := fmt.Sprintf("%s\n%s\nЗП: %d-%d\n%s",
				vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)
			h.bot.SendMessage(userId, msg)
		}
	}

}

func (h *Handler) handleUnknown(userId int) {
	log.Printf("unknown command")
	h.bot.SendMessage(userId, "Ошибка при отправке сообщения")
}
