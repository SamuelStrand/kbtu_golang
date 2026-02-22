package modules

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() AppConfig {
	_ = godotenv.Load()

	serverAddr := getEnv("SERVER_ADDR", ":8080")
	apiKey := getEnv("API_KEY", "dev-api-key")

	execTimeout, err := time.ParseDuration(getEnv("PG_TIMEOUT", "5s"))
	if err != nil {
		execTimeout = 5 * time.Second
	}

	pg := PostgreConfig{
		Host:        getEnv("PG_HOST", "localhost"),
		Port:        getEnv("PG_PORT", "5432"),
		Username:    getEnv("PG_USER", "postgres"),
		Password:    getEnv("PG_PASSWORD", "postgres"),
		DBName:      getEnv("PG_DB", "mydb"),
		SSLMode:     getEnv("PG_SSLMODE", "disable"),
		ExecTimeout: execTimeout,
	}

	return AppConfig{
		ServerAddr: serverAddr,
		APIKey:     apiKey,
		Postgres:   pg,
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
