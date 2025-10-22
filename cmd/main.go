package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	apppkg "tgbot/internal/app"
	cfgpkg "tgbot/internal/config"
	"tgbot/internal/logger"
)

func main() {
	cfg := cfgpkg.LoadConfig()
	slogger := logger.GetLogger(cfg.ENV)
	slogger.Info("config and logger loaded")

	// build app (app.New собирает все зависимости и возвращает runnable App)
	app, err := apppkg.New(cfg, slogger)
	if err != nil {
		slogger.Error("failed to init app:", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run returns when context is cancelled or fatal error occurs
	runErr := make(chan error, 1)
	go func() {
		runErr <- app.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		// graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = app.Shutdown(shutdownCtx)
	case err := <-runErr:
		// fatal error from app
		_ = app.Shutdown(context.Background())
		slogger.Error("app run failed", "error", err)
	}
}
