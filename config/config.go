package config

import (
	"os"
)

type Config struct {
	Port      string
	TaxAPIURL string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		TaxAPIURL: getEnv("TAX_API_URL", "http://tax-api:5001"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
