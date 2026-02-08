package main

import (
	"log"
	"net/http"
	"time"

	"practice2/internal/handlers"
	"practice2/internal/middleware"
	"practice2/internal/store"
)

func main() {
	taskStore := store.NewTaskStore()

	mux := http.NewServeMux()
	mux.Handle("/tasks", handlers.NewTaskHandler(taskStore))

	var h http.Handler = mux
	h = middleware.Auth(h, "secret12345")
	h = middleware.Logging(h, "Task API request")

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server started on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
