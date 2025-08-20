package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	if os.Getenv("APP_PORT") == "" {
		log.Println("No .env file found or variables missing; using environment/defaults")
	}

	return &Config{
		AppPort: first(os.Getenv("APP_PORT"), "8080"),
		DBHost:  first(os.Getenv("DB_HOST"), "localhost"),
		DBPort:  first(os.Getenv("DB_PORT"), "5432"),
		DBUser:  first(os.Getenv("DB_USER"), "postgres"),
		DBPass:  first(os.Getenv("DB_PASSWORD"), "postgres"),
		DBName:  first(os.Getenv("DB_NAME"), "subscriptions"),
	}
}

func first(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
