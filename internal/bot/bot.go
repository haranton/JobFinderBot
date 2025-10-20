package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const telegramAPI = "https://api.telegram.org/bot"

type BotClient interface {
	SendMessage(chatID int64, text string) error
	SendChatAction(chatID int64, action string) error
}

type Bot struct {
	Token string
}

func (b *Bot) SendMessage(chatID int, text string) error {
	params := url.Values{}
	params.Add("chat_id", fmt.Sprintf("%d", chatID))
	params.Add("text", text)

	_, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", telegramAPI, b.Token), params)
	return err
}

func (b *Bot) RegisterCommands() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setMyCommands", b.Token)

	type Command struct {
		Command     string `json:"command"`
		Description string `json:"description"`
	}

	commands := []Command{
		{Command: "start", Description: "Регистрация пользователя"},
		{Command: "help", Description: "Список команд"},
		{Command: "deletesubscribes", Description: "Удалить подписки"},
		{Command: "subscribes", Description: "Посмотреть текущие подписки"},
	}

	body, err := json.Marshal(map[string]interface{}{
		"commands": commands,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("telegram API returned status %d: %s", resp.StatusCode, buf.String())
	}

	return nil
}
