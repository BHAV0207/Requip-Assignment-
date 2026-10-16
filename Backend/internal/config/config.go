package config

import (
	"log"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	LogLevel string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	DBMaxOpenConns int
	DBMaxIdleConns int

	// RunMigrations controls whether the API runs GORM AutoMigrate on startup.
	RunMigrations bool

	CORSAllowedOrigins []string

	// FaultErrorRate makes this fraction of /api requests fail with 500.
	// Used only to demonstrate canary analysis and automatic rollback.
	FaultErrorRate float64
}

func Load() *Config {

	// .env is a local-development convenience. In Kubernetes the values come
	// from a ConfigMap and a Secret, so a missing file is expected there.
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}

	return &Config{
		AppPort:  getEnv("APP_PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "user_management"),

		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),

		RunMigrations: getEnv("RUN_MIGRATIONS", "true") == "true",

		CORSAllowedOrigins: splitList(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),

		FaultErrorRate: getEnvFloat("FAULT_ERROR_RATE", 0),
	}
}

func getEnvInt(key string, fallback int) int {
	v, err := strconv.Atoi(getEnv(key, ""))
	if err != nil {
		return fallback
	}
	return v
}

func getEnvFloat(key string, fallback float64) float64 {
	v, err := strconv.ParseFloat(getEnv(key, ""), 64)
	if err != nil {
		return fallback
	}
	return v
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
