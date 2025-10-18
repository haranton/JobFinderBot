package bot

import (
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
	Token  string
	client *http.Client
}

func (b *Bot) SendMessage(chatID int, text string) error {
	params := url.Values{}
	params.Add("chat_id", fmt.Sprintf("%d", chatID))
	params.Add("text", text)

	_, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", telegramAPI, b.Token), params)
	return err
}
