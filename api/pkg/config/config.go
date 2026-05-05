package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	AppEnv   string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	RedisAddr string
	RedisPass string
}

func Load() *Config {
	godotenv.Load("../.env")
	return &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		AppEnv:    getEnv("APP_ENV", "development"),
		DBHost:    getEnv("POSTGRES_HOST", "localhost"),
		DBPort:    getEnv("POSTGRES_PORT", "5432"),
		DBUser:    getEnv("POSTGRES_USER", "sentinel"),
		DBPass:    getEnv("POSTGRES_PASSWORD", "changeme"),
		DBName:    getEnv("POSTGRES_DB", "sentinel20"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASSWORD", "changeme"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
