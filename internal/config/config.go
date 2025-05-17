package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Environment string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}