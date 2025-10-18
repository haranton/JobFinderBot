package config

import (
	"log"
	"os"
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
	REDIS_HOST      string
	REDIS_PORT      string
	REDIS_PASSWORD  string
	REDIS_DB        int
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &Config{
		AppPort:         os.Getenv("APP_PORT"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		ENV:             os.Getenv("ENV"),
		MIGRATIONS_PATH: os.Getenv("MIGRATIONS_PATH"),
		REDIS_HOST:      os.Getenv("REDIS_HOST"),
		REDIS_PORT:      os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD:  os.Getenv("REDIS_PASSWORD"),
		REDIS_DB:        redisDb,
	}
}
