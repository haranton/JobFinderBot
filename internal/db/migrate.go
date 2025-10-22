package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"tgbot/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Config, logger *slog.Logger) error {

	DBHost := cfg.DBHost
	if cfg.ENV == "PRODUCTION" {
		DBHost = cfg.DBHostProd
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBPort == "" || cfg.DBName == "" {
		return fmt.Errorf("incomplete DB configuration: user=%q, name=%q, host=%q",
			cfg.DBUser, cfg.DBName, DBHost)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, DBHost, cfg.DBPort, cfg.DBName,
	)

	logger.Info("Starting database migrations")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			logger.Warn("failed to close DB connection", slog.String("error", cerr.Error()))
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}

	// Определяем путь к миграциям
	migrationPath := cfg.MIGRATIONS_PATH
	if cfg.ENV == "PRODUCTION" {
		migrationPath = "/app/migrations"
	} else {
		path, err := findMigrationsPath()
		if err != nil {
			return fmt.Errorf("find migrations: %w", err)
		}
		migrationPath = path
	}

	logger.Info("Using migrations path", slog.String("path", migrationPath))

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No new migrations to apply")
		} else {
			return fmt.Errorf("run migrations: %w", err)
		}
	}

	logger.Info("Database migrations ran successfully")
	return nil
}

func findMigrationsPath() (string, error) {
	ex, _ := os.Executable()
	basePath := filepath.Dir(ex)
	possiblePaths := []string{
		filepath.Join(basePath, "migrations"),
		filepath.Join(basePath, "..", "migrations"),
		filepath.Join(".", "migrations"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("migrations folder not found (searched from %s)", basePath)
}
