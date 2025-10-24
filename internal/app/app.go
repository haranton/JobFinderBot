package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"tgbot/internal/bot"
	botTg "tgbot/internal/bot"
	"tgbot/internal/config"
	"tgbot/internal/db"
	"tgbot/internal/dto"
	"tgbot/internal/fetcher"
	"tgbot/internal/handler"
	"tgbot/internal/repo"
	"tgbot/internal/sender"
	"tgbot/internal/service"
)

type App struct {
	cfg        *config.Config
	logger     *slog.Logger
	repo       *repo.Repository
	bot        *botTg.Bot
	svc        *service.Service
	sender     *sender.Sender
	fetcher    *fetcher.Fetcher
	handler    *handler.Handler
	httpClient *http.Client
	wg         sync.WaitGroup
}

func New(cfg *config.Config, slogger *slog.Logger) (*App, error) {
	dbConn := db.GetDBConnect(cfg, slogger)
	if err := db.RunMigrations(cfg, slogger); err != nil {
		return nil, err
	}

	rep := repo.NewRepository(dbConn)
	f := fetcher.NewFetcher()
	b := botTg.NewBot(cfg.TOKEN, slogger)
	if err := b.RegisterCommands(); err != nil {
		return nil, err
	}

	svc := service.NewService(rep, f)
	s := sender.NewSender(svc, b, slogger)
	hand := handler.NewHandler(svc, b, slogger)

	return &App{
		cfg:        cfg,
		logger:     slogger,
		repo:       rep,
		bot:        b,
		svc:        svc,
		sender:     s,
		fetcher:    f,
		httpClient: &http.Client{Timeout: 35 * time.Second},
		handler:    hand,
	}, nil
}

// Run запускает poller, sender и worker pool и возвращает, когда ctx отменён или произойдёт фатал ошибка.
func (a *App) Run(ctx context.Context) error {
	jobs := make(chan dto.Update, 50)

	// start sender (предполагается, что sender.Start(ctx) корректно слушает ctx)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.sender.Start(ctx)
		a.logger.Info("sender start")
	}()

	// start workers
	for i := 0; i < a.cfg.WorkerCount; i++ {
		a.wg.Add(1)
		go func(workerID int) {
			defer a.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					ctxTime, cancel := context.WithTimeout(ctx, 10*time.Second)
					a.handler.HandleMessage(ctxTime, job.Message.Chat.ID, job.Message.Text)
					cancel()
				}
			}
		}(i)
	}

	// start poller
	poller := bot.NewPoller(a.httpClient, a.cfg.TOKEN)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := poller.Start(ctx, jobs); err != nil {
			a.logger.Error("poller stopped", "error", err)
		}
	}()

	// wait until ctx cancelled
	<-ctx.Done()
	// при остановке закрываем канал jobs чтобы воркеры корректно завершились
	close(jobs)

	// ждём завершения goroutines
	a.wg.Wait()
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application...")
	done := make(chan struct{})

	go func() {
		fmt.Println("жлем 3 секунды")
		time.Sleep(3 * time.Second)
		a.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("timeout for shutdown")
		return ctx.Err()
	case <-done:
		a.logger.Info("shutdown success")
	}

	if err := a.repo.Close(); err != nil {
		a.logger.Error("Failed to close database", "error", err)
		return err
	}
	a.logger.Info("Database connection closed successfully")
	return nil
}
