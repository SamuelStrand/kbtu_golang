package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"practice5/internal/httpapi"
	"practice5/internal/repository"
	"practice5/internal/usecase"
	"practice5/pkg/config"
)

type App struct {
	Router http.Handler
}

func New(cfg config.Config) (*App, error) {
	db, err := connectWithRetry(cfg.PostgresDSN(), 25, 600*time.Millisecond)
	if err != nil {
		return nil, err
	}

	if err := repository.EnsureSchema(context.Background(), db); err != nil {
		return nil, err
	}
	if err := repository.SeedIfNeeded(context.Background(), db); err != nil {
		return nil, err
	}

	repo := repository.NewPostgresUserRepo(db)
	uc := usecase.NewUsersUsecase(repo)
	h := httpapi.NewHandler(uc)

	mux := httpapi.NewRouter(h)

	return &App{Router: mux}, nil
}

func connectWithRetry(dsn string, attempts int, delay time.Duration) (*sql.DB, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			lastErr = err
			log.Printf("db open failed: %v", err)
			time.Sleep(delay)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			return db, nil
		}
		_ = db.Close()
		lastErr = err
		log.Printf("waiting for database... (%v)", err)
		time.Sleep(delay)
	}
	return nil, lastErr
}
