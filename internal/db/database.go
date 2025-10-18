package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"tgbot/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func GetDBConnect(config *config.Config, logger *slog.Logger) *sqlx.DB {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	logger.Info("Connecting to database", slog.String("dsn", dsn))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect database")
		os.Exit(1)
	}

	return db

}

func GetRedisConnect(config *config.Config, logger *slog.Logger) *redis.Client {

	ctx := context.Background()
	dsn := fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT)

	logger.Info("Connecting to Redis", slog.String("dsn", dsn))

	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,                   // адрес Redis
		Password: config.REDIS_PASSWORD, // пароль (если есть)
		DB:       config.REDIS_DB,       // номер базы (по умолчанию 0)
	})

	// проверяем соединение
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error("failed to connect to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Connected to Redis successfully")
	return rdb
}
