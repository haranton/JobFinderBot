package db

import (
	"fmt"
	"log/slog"
	"os"
	"tgbot/internal/config"

	"github.com/jmoiron/sqlx"
)

func GetDBConnect(config *config.Config, logger *slog.Logger) *sqlx.DB {

	var DBHost string
	if config.ENV == "PRODUCTION" {
		DBHost = config.DBHostProd
	} else {
		DBHost = config.DBHost
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost,
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
