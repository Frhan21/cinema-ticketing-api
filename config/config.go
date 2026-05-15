package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env can't be loaded")
	}
}

func GetEnv(key string, defaultValue string) string {

	value := os.Getenv(key)

	if value == "" {

		return defaultValue
	}
	return value
}
