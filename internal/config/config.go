package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	return &Config{
		Port:        getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Validate проверяет обязательные настройки приложения.
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("переменная DATABASE_URL не задана")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("переменная JWT_SECRET не задана")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET должен содержать минимум 32 символа")
	}
	return nil
}
