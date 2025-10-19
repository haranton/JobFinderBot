package config

import (
	"log"
	"os"
	"path/filepath"

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

	return &Config{
		AppPort:         os.Getenv("APP_PORT"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		ENV:             os.Getenv("ENV"),
		MIGRATIONS_PATH: os.Getenv("MIGRATIONS_PATH"),
		TOKEN:           os.Getenv("TOKEN"),
	}
}
