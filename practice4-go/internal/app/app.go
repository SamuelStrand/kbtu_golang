package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	httpHandler "practice4/internal/handler/http"
	"practice4/internal/middleware"
	"practice4/internal/repository"
	"practice4/internal/repository/_postgres"
	usersUC "practice4/internal/usecase/users"
	"practice4/pkg/modules"
)

func Run() {
	cfg := modules.LoadConfig()
	logger := log.New(os.Stdout, "", 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := _postgres.NewPGXDialect(ctx, &cfg.Postgres)
	defer func() { _ = db.DB.Close() }()

	repos := repository.NewRepositories(db)
	usersUsecase := usersUC.New(repos.UserRepository)
	h := httpHandler.New(usersUsecase)

	r := chi.NewRouter()
	r.Use(middleware.Logging(logger))

	r.Get("/health", h.Health)
	r.Get("/swagger", h.SwaggerUI)
	r.Get("/swagger/swagger.json", h.SwaggerSpec)

	r.Group(func(pr chi.Router) {
		pr.Use(middleware.APIKey(cfg.APIKey))

		pr.Route("/users", func(r chi.Router) {
			r.Get("/", h.GetUsers)
			r.Post("/", h.CreateUser)
			r.Get("/{id}", h.GetUserByID)
			r.Put("/{id}", h.UpdateUser)
			r.Patch("/{id}", h.PatchUser)
			r.Delete("/{id}", h.DeleteUser)
		})
	})

	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Printf("%s %s %s", time.Now().Format(time.RFC3339), "START", cfg.ServerAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}
