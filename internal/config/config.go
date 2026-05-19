package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig
	DB   DBConfig
	JWT  JWTConfig
	SMTP SMTPConfig
	Admin AdminConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Sender   string
}

type AdminConfig struct {
	Name     string
	Email    string
	Password string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", ":8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "cinema_ticketing"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "secret_key"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnvInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			Sender:   getEnv("SMTP_SENDER", ""),
		},
		Admin: AdminConfig{
			Name:     getEnv("ADMIN_NAME", "Administrator"),
			Email:    getEnv("ADMIN_EMAIL", ""),
			Password: getEnv("ADMIN_PASSWORD", ""),
		},
	}
}

// getEnv is an internal helper to read env with a fallback default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt is an internal helper to read env as integer with a fallback default value.
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: invalid integer value for %s, using default %d\n", key, defaultValue)
		return defaultValue
	}
	return intValue
}
