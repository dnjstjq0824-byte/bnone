package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	DBDriver     string // "postgres" or "sqlite"
	JWTSecret    string
	Environment  string
	BaseCurrency string
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "bn_one.db")
	driver := getEnv("DB_DRIVER", "sqlite")
	if len(dbURL) > 10 && dbURL[:8] == "postgres" {
		driver = "postgres"
	}

	return &Config{
		Port:         port,
		DatabaseURL:  dbURL,
		DBDriver:     driver,
		JWTSecret:    getEnv("JWT_SECRET", "bn-one-super-secret-jwt-key-2026-production-grade"),
		Environment:  getEnv("APP_ENV", "development"),
		BaseCurrency: getEnv("BASE_CURRENCY", "KRW"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
