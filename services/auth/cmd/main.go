package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NiklasWillecke/media-platform/services/auth/internal/handler"
	db "github.com/niklaswillecke/streaming-platform/shared/db"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	dbConn := db.NewDB("postgres://postgres:example@localhost:5432/postgres?sslmode=disable")
	defer dbConn.Close()

	handler := handler.NewHandler(dbConn.Q)
	mux := http.NewServeMux()
	handler.RegisterHandler(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server is starting...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v\n", err)
	}

	log.Println("Shutdown completed.")

}
