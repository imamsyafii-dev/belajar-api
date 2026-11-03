package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL wajib diisi")
	}

	pool, err := NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repository := NewPostgresUserRepository(pool)
	service := NewUserService(repository)
	userHandler := NewUserHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.ListHandler)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByIDHandler)
	mux.HandleFunc("POST /users", userHandler.CreateHandler)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateHandler)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Println("server berjalan di http://localhost:8080")

		serverError <- server.ListenAndServe()
	}()

	shutdownSignal, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}

	case <-shutdownSignal.Done():
		log.Println("shutdown signal diterima")

		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			log.Printf("graceful shutdown gagal: %v", err)
		}

		log.Println("server berhenti dengan aman")
	}
}
