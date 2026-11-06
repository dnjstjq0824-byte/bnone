package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	DBDriver               string // "postgres" or "sqlite"
	JWTSecret              string
	Environment            string
	BaseCurrency           string
	LogLevel               string
	CORSAllowedOrigins     []string
	AdminBootstrapEnabled  bool
	AdminBootstrapPassword string
}

// devOnlyJWTSecret is a development fallback. Production refuses to boot
// without an explicit JWT_SECRET (see Load).
const devOnlyJWTSecret = "moa-dev-only-jwt-key-do-not-use-in-production"

func Load() *Config {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "moa.db")
	driver := getEnv("DB_DRIVER", "sqlite")
	if len(dbURL) > 10 && dbURL[:8] == "postgres" {
		driver = "postgres"
	}

	// .env.example documents ENVIRONMENT; accept the legacy APP_ENV alias too.
	env := getEnv("ENVIRONMENT", "")
	if env == "" {
		env = getEnv("APP_ENV", "development")
	}
	isProd := env == "production" || env == "prod"

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if isProd {
			panic("MOA: JWT_SECRET must be set explicitly in production (min 32 chars). Refusing to boot with a default key.")
		}
		jwtSecret = devOnlyJWTSecret
	}
	if isProd && len(jwtSecret) < 32 {
		panic("MOA: JWT_SECRET must be at least 32 characters in production.")
	}

	return &Config{
		Port:                   port,
		DatabaseURL:            dbURL,
		DBDriver:               driver,
		JWTSecret:              jwtSecret,
		Environment:            env,
		BaseCurrency:           getEnv("BASE_CURRENCY", "KRW"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		CORSAllowedOrigins:     splitOrigins(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")),
		AdminBootstrapEnabled:  getEnv("ADMIN_BOOTSTRAP_ENABLED", "false") == "true",
		AdminBootstrapPassword: os.Getenv("ADMIN_BOOTSTRAP_PASSWORD"),
	}
}

func splitOrigins(raw string) []string {
	var out []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" || part == "*" {
			// Never allow a wildcard origin: it cannot be combined
			// with credentialed requests.
			continue
		}
		out = append(out, part)
	}
	return out
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

