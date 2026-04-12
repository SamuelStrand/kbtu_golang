package app

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"practice7/internal/httpserver"
	"practice7/internal/middleware"
	"practice7/internal/store"
	"practice7/internal/utils"

	"github.com/gin-gonic/gin"
)

type Config struct {
	HTTPPort          string
	JWTSecret         string
	RateLimit         int
	RateWindowSeconds int
}

func loadConfig() Config {
	return Config{
		HTTPPort:          getenv("HTTP_PORT", "8080"),
		JWTSecret:         getenv("JWT_SECRET", "dev-secret-change-me"),
		RateLimit:         getenvInt("RATE_LIMIT", 5),
		RateWindowSeconds: getenvInt("RATE_WINDOW_SECONDS", 10),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func Run() error {
	cfg := loadConfig()
	secret := []byte(cfg.JWTSecret)

	userStore := store.NewInMemoryUserStore()
	auth := utils.NewAuthService(userStore, secret)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Global rate limiter (applies to all endpoints)
	rl := middleware.NewRateLimiter(cfg.RateLimit, cfg.RateWindowSeconds, secret)
	r.Use(rl.Middleware())

	httpserver.RegisterRoutes(r, auth, secret)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("server listening on %s", addr)
	return r.Run(addr)
}
