package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"tgbot/internal/bot"
	"tgbot/internal/config"
	"tgbot/internal/db"
	"tgbot/internal/fetcher"
	"tgbot/internal/handler"
	"tgbot/internal/logger"
	"tgbot/internal/repo"
	"tgbot/internal/sender"
	"tgbot/internal/service"
	"time"
)

const telegramAPI = "https://api.telegram.org/bot"

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
		From struct {
			Username string `json:"username"`
		} `json:"from"`
	} `json:"message"`
}

func main() {
	//config
	config := config.LoadConfig()
	//Logger
	logger := logger.GetLogger(config.ENV)

	logger.Info("config and logger successfully load")
	//Db
	ConnectDb := db.GetDBConnect(config, logger)
	defer ConnectDb.Close()

	// migrations
	db.RunMigrations(config, logger)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	fetcher := fetcher.NewFetcher(client)

	bot := &bot.Bot{Token: config.TOKEN}

	if err := bot.RegisterCommands(); err != nil {
		log.Fatalf("failed to register bot commands: %v", err)
	}

	repo := repo.NewRepository(ConnectDb)
	service := service.NewService(repo, fetcher)
	handler := handler.NewHandler(service, bot)

	sender := sender.NewSender(service, bot)
	sender.Start()

	offset := 0

	for {
		updates, err := getUpdates(config.TOKEN, offset)
		if err != nil {
			log.Println("error get updates")
			continue
		}

		for _, update := range updates {
			id := update.UpdateID
			text := update.Message.Text
			username := update.Message.From.Username
			chatID := update.Message.Chat.ID

			handler.HandleMessage(chatID, text)
			log.Printf("[%s] %s номер сообщения-%v", username, text, id)

			offset = update.UpdateID + 1
		}
	}
}

func getUpdates(token string, offset int) ([]Update, error) {

	resp, err := http.Get(fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=30", telegramAPI, token, offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Result, nil

}
