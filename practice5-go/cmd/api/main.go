// @title   Practice5 API
// @version 1.0
// @host    localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"
	"time"

	_ "practice5/docs"

	"practice5/internal/app"
	"practice5/pkg/config"
)

func main() {
	config.LoadDotEnv(".env")

	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           application.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server started on http://localhost:%s", cfg.HTTPPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
