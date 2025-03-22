package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}

func GetConfig() *Config {
	godotenv.Load()

	return &Config{
		HOST:     GetEnv("DB_HOST", "localhost"),
		PORT:     GetEnv("DB_PORT", "5432"),
		USER:     GetEnv("DB_USER", "postgres"),
		PASSWORD: GetEnv("DB_PASSWORD", "password"),
		DBNAME:   GetEnv("DB_NAME", "postgres"),
	}
}

func GetEnv(variableName string, defaultValue string) string {
	val, exists := os.LookupEnv(variableName)
	if exists {
		return val
	}
	return defaultValue
}
