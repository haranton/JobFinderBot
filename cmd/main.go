package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	fetcher := fetcher.NewFetcher()

	bot := &bot.Bot{Token: config.TOKEN}

	if err := bot.RegisterCommands(); err != nil {
		log.Fatalf("failed to register bot commands: %v", err)
	}

	repo := repo.NewRepository(ConnectDb)
	service := service.NewService(repo, fetcher)
	handler := handler.NewHandler(service, bot)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sender := sender.NewSender(service, bot)
	sender.Start(ctx)

	offset := 0
	jobs := make(chan Update, 50)

	for i := 0; i < config.WorkerCount; i++ {
		go Worker(jobs, handler)
	}

	go func() {
		for {
			updates, err := getUpdates(config.TOKEN, offset)
			if err != nil {
				log.Println("error get updates")
				continue
			}

			for _, update := range updates {
				jobs <- update
				offset = update.UpdateID + 1
			}
		}
	}()

	<-stop
	cancel()
	close(jobs)
	time.Sleep(4 * time.Second)
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

func Worker(jobs chan Update, handler *handler.Handler) {

	for job := range jobs {

		id := job.UpdateID
		text := job.Message.Text
		username := job.Message.From.Username
		chatID := job.Message.Chat.ID

		handler.HandleMessage(chatID, text)
		log.Printf("[%s] %s номер сообщения-%v", username, text, id)
	}
}
