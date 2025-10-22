package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"tgbot/internal/bot"
	"tgbot/internal/models"
	"tgbot/internal/service"
)

type Handler struct {
	bot     *bot.Bot
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(service *service.Service, bot *bot.Bot, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		bot:     bot,
		logger:  logger,
	}
}

// HandleMessage принимает общий контекст от воркера
func (h *Handler) HandleMessage(ctx context.Context, userId int, text string) {
	text = strings.TrimSpace(text)

	switch {
	case text == "/start":
		h.handleStart(ctx, userId)
	case strings.HasPrefix(text, "/find"):
		h.handleFind(ctx, userId, text)
	case text == "/subscribes":
		h.handleSubscribes(ctx, userId)
	case strings.HasPrefix(text, "/subscribe "):
		h.handleSubscribe(ctx, userId, text)
	case text == "/deletesubscribes":
		h.handleDeleteSubscribes(ctx, userId)
	case text == "/help":
		h.handleHelp(userId)
	default:
		h.handleUnknown(userId)
	}
}

// ---------- USERS ----------

func (h *Handler) handleStart(ctx context.Context, userId int) {
	user, err := h.service.RegisterUser(ctx, userId)
	if err != nil {
		h.logger.Error("Ошибка регистрации пользователя", "userId", userId, "error", err)

		// если пользователь уже есть — отвечаем корректно
		if user != (models.User{}) {
			h.bot.SendMessage(userId, "Пользователь уже зарегистрирован")
			return
		}

		h.bot.SendMessage(userId, "Ошибка регистрации пользователя")
		return
	}

	msg := fmt.Sprintf("Пользователь успешно зарегистрирован: %v", user)
	h.bot.SendMessage(userId, msg)
}

// ---------- SUBSCRIPTIONS ----------

func (h *Handler) handleDeleteSubscribes(ctx context.Context, userId int) {
	if err := h.service.DeleteSubscribes(ctx, userId); err != nil {
		h.logger.Error("Ошибка при удалении подписок", "userId", userId, "error", err)
		h.bot.SendMessage(userId, "ошибка при удалении подписок")
		return
	}
	h.bot.SendMessage(userId, "все подписки удалены успешно")
}

func (h *Handler) handleSubscribes(ctx context.Context, userId int) {
	subscribes, err := h.service.SubscriptionsUser(ctx, userId)
	if err != nil {
		h.logger.Error("Ошибка при получении подписок", "userId", userId, "error", err)
		h.bot.SendMessage(userId, "ошибка при получении подписок")
		return
	}

	if len(subscribes) == 0 {
		h.bot.SendMessage(userId, "Подписки отсутствуют")
		return
	}

	var msg strings.Builder
	for _, sub := range subscribes {
		msg.WriteString(fmt.Sprintf("Подписка: %s\n", sub.SearchText))
	}

	h.bot.SendMessage(userId, msg.String())
}

func (h *Handler) handleSubscribe(ctx context.Context, userId int, subscribeCommandText string) {
	subscribeText := strings.TrimPrefix(subscribeCommandText, "/subscribe ")

	if strings.TrimSpace(subscribeText) == "" {
		h.bot.SendMessage(userId, "не введены ключевые слова для подписки")
		return
	}

	subscribe, err := h.service.RegisterSubscribe(ctx, userId, subscribeText)
	if err != nil {
		h.logger.Error("Ошибка регистрации подписки", "userId", userId, "error", err)
		h.bot.SendMessage(userId, "Ошибка регистрации подписки")
		return
	}

	msg := fmt.Sprintf("Подписка успешно зарегистрирована: %v", subscribe)
	h.bot.SendMessage(userId, msg)
}

// ---------- VACANCIES ----------

func (h *Handler) handleFind(ctx context.Context, userId int, text string) {
	query := strings.TrimPrefix(text, "/find ")

	if strings.TrimSpace(query) == "" {
		h.bot.SendMessage(userId, "не введены ключевые слова для поиска")
		return
	}

	vacancies, err := h.service.SearchVacancies(ctx, query, userId)
	if err != nil {
		h.logger.Error("Ошибка при получении вакансий", "query", query, "error", err)
		h.bot.SendMessage(userId, "Ошибка при поиске вакансий")
		return
	}

	for _, vac := range vacancies {
		msg := fmt.Sprintf("%s\n%s\nЗП: %d-%d\n%s",
			vac.Name, vac.Area.Name, vac.Salary.From, vac.Salary.To, vac.Url)
		h.bot.SendMessage(userId, msg)
	}
}

// ---------- HELP / UNKNOWN ----------

func (h *Handler) handleUnknown(userId int) {
	h.logger.Warn("Неизвестная команда", "userId", userId)
	h.bot.SendMessage(userId, "неизвестная команда, используйте /help")
}

func (h *Handler) handleHelp(userId int) {
	helpText := `*Список доступных команд:*

	/start — регистрация пользователя  
	/find <запрос> — поиск вакансий  
	/subscribe <запрос> — подписка на вакансии  
	/subscribes — список подписок  
	/deletesubscribes — удалить все подписки  

	*Примеры:*  
	/find golang Ижевск  
	/subscribe python developer удаленно`

	h.bot.SendMessage(userId, helpText)
}
