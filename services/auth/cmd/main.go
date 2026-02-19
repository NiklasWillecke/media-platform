package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "streaming-platform/shared/db"

	"streaming-platform/services/auth/internal/handler"
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

	err := dbConn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Cant connect to DB: %v\n", err)
	}

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
