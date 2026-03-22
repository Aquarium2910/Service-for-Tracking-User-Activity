package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DbHost         string
	DbPort         string
	DbUser         string
	DbPassword     string
	DbName         string
	WorkerInterval time.Duration
}

func LoadConfig(logger *slog.Logger) *Config {
	_ = godotenv.Load()

	return &Config{
		Port:           requireEnv("PORT", logger),
		DbHost:         requireEnv("DB_HOST", logger),
		DbPort:         requireEnv("DB_PORT", logger),
		DbUser:         requireEnv("DB_USER", logger),
		DbPassword:     requireEnv("DB_PASSWORD", logger),
		DbName:         requireEnv("DB_NAME", logger),
		WorkerInterval: requireDurationEnv("WORKER_INTERVAL", logger),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
}

func requireEnv(key string, logger *slog.Logger) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		logger.Error("Required environment variable is missing or empty", slog.String("key", key))
		os.Exit(1)
	}
	return value
}

func requireDurationEnv(key string, logger *slog.Logger) time.Duration {
	strValue := requireEnv(key, logger)
	duration, err := time.ParseDuration(strValue)
	if err != nil {
		logger.Error("Invalid format for duration environment variable", slog.String("key", key),
			slog.String("error", err.Error()))
		os.Exit(1)
	}
	return duration
}
