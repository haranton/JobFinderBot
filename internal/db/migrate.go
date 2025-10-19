package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"tgbot/internal/config"

	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Config, logger *slog.Logger) {
	// Определяем хост базы данных
	var DBHost string
	if cfg.ENV == "PRODUCTION" {
		DBHost = cfg.DBHostProd
	} else {
		DBHost = cfg.DBHost
	}

	// Формируем DSN строку
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, DBHost, cfg.DBPort, cfg.DBName,
	)

	logger.Info("Starting database migrations")

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to open database connection", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("failed to create migration driver", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Определяем путь к миграциям
	migrationPath := cfg.MIGRATIONS_PATH
	if cfg.ENV == "PRODUCTION" {
		// В контейнере миграции лежат в /app/migrations
		migrationPath = "/app/migrations"
	} else if migrationPath == "" {
		// Для локального запуска
		migrationPath = "./migrations"
	}

	// Проверяем, существует ли папка
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		// Попробуем вычислить путь относительно исполняемого файла
		ex, _ := os.Executable()
		basePath := filepath.Dir(ex)
		altPath := filepath.Join(basePath, "migrations")
		if _, err2 := os.Stat(altPath); err2 == nil {
			migrationPath = altPath
		} else {
			logger.Error("migrations folder not found", slog.String("checked_path", migrationPath))
			os.Exit(1)
		}
	}

	logger.Info("Using migrations path", slog.String("path", migrationPath))

	// Создаём экземпляр мигратора
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver,
	)
	if err != nil {
		logger.Error("failed to create migrate instance", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Запускаем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Database migrations ran successfully")
}
