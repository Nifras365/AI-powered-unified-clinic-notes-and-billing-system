package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	Port                 string
	DatabaseURL          string
	GeminiAPIKey         string
	GeminiModel          string
	BaseConsultationFee  float64
	DefaultDiscountType  string
	DefaultDiscountValue float64
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		GeminiAPIKey:         os.Getenv("GEMINI_API_KEY"),
		GeminiModel:          getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		BaseConsultationFee:  getEnvAsFloat("BASE_CONSULTATION_FEE", 40),
		DefaultDiscountType:  getEnv("DEFAULT_DISCOUNT_TYPE", "fixed"),
		DefaultDiscountValue: getEnvAsFloat("DEFAULT_DISCOUNT_VALUE", 0),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsFloat(key string, fallback float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fallback
	}

	return value
}
