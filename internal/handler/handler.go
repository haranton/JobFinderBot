package handler

import (
	"fmt"
	"log"
	"strings"
	"tgbot/internal/bot"
	"tgbot/internal/models"
	"tgbot/internal/service"
)

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
	case text == "/subscribes":
		h.handleSubscribes(userId)
	case strings.HasPrefix(text, "/subscribe "):
		h.handleSubscribe(userId, text)
	case text == "/deletesubscribes":
		h.handleDeleteSubscribes(userId)
	case text == "/help":
		h.handleHelp(userId)
	default:
		h.handleUnknown(userId)
	}
}

func (h *Handler) handleDeleteSubscribes(userId int) {
	err := h.service.DeleteSubscribes(userId)
	if err != nil {
		h.bot.SendMessage(userId, "ошибка при удалении подписок")
		return
	}

	h.bot.SendMessage(userId, "все подписки удалены успешно")

}

func (h *Handler) handleSubscribes(userId int) {
	subscribes, err := h.service.SubscriptionsUser(userId)
	if err != nil {
		h.bot.SendMessage(userId, "ошибка при получении подписок")
		return
	}

	if len(subscribes) == 0 {
		h.bot.SendMessage(userId, "Подписки отсутствуют")
		return
	}

	var msg string
	for _, sub := range subscribes {
		msg = msg + fmt.Sprintf("Подписка: %s\n", sub.SearchText)
	}

	h.bot.SendMessage(userId, msg)

}

func (h *Handler) handleSubscribe(userId int, subscribeCommandText string) {

	subscribeText := strings.TrimPrefix(subscribeCommandText, "/subscribe")

	if subscribeText == "" {
		h.bot.SendMessage(userId, "не введены ключевые слова для подписки")
		return
	}
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

	if query == "" {
		h.bot.SendMessage(userId, "не введены ключевые слова для подписки")
		return
	}

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
	h.bot.SendMessage(userId, "не знакомая команда, посмотрите список команд /help")
}

func (h *Handler) handleHelp(userId int) {
	helpText := `*Список доступных команд:*

		/start — регистрация пользователя  
		/find <запрос> — поиск вакансий по ключевым словам  
		/subscribe <запрос> — подписка на вакансии  

		*Примеры:*
		/find golang Ижевск
		/subscribe python developer удаленно
		`

	if err := h.bot.SendMessage(userId, helpText); err != nil {
		log.Println("error sending help message:", err)
	}
}
