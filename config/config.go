package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Environment string
}

func LoadConfig() *Config {
	// Load .env file if exists (optional)
	godotenv.Load()

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=myuser password=mypassword dbname=usersdb port=5432 sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		Port:        getEnv("PORT", "3001"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
