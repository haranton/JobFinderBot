package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	ENV             string
	MIGRATIONS_PATH string
	TOKEN           string
	DBHostProd      string
	WorkerCount     int
}

func LoadConfig() *Config {

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)

	// Возможные пути к .env
	paths := []string{
		filepath.Join(exPath, "../../.env"), // запуск из /cmd
		filepath.Join(exPath, "../.env"),    // запуск из корня
		".env",                              // fallback
	}

	loaded := false
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			err = godotenv.Load(p)
			if err != nil {
				log.Fatalf("Ошибка загрузки .env: %v", err)
			}
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println(" .env файл не найден — использую переменные окружения из системы")
	}

	workerCountString := os.Getenv("WORKER_COUNT")
	workerCount, err := strconv.Atoi(workerCountString)
	if err != nil {
		log.Fatalf("failed to parce workerCount in integer, err: %w", err)
	}

	return &Config{
		AppPort:         os.Getenv("APP_PORT"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBHostProd:      os.Getenv("DB_HOST_PROD"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		ENV:             os.Getenv("ENV"),
		MIGRATIONS_PATH: os.Getenv("MIGRATIONS_PATH"),
		TOKEN:           os.Getenv("TOKEN"),
		WorkerCount:     workerCount,
	}
}
