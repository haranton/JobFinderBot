package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"tgbot/internal/config"

	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Config, logger *slog.Logger) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	logger.Info("Starting database migrations")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to open database connection:", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("failed to create migration driver:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	basepath := cfg.MIGRATIONS_PATH
	if basepath == "" {
		logger.Error("MIGRATIONS_PATH is not set in the configuration")
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+basepath,
		"postgres", driver,
	)
	if err != nil {
		logger.Error("failed to create migrate instance:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to run migrations:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Database migrations ran successfully")
}
