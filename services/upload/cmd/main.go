package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"streaming-platform/services/upload/pkg/handler"
	"streaming-platform/services/upload/pkg/queue"
	dataStore "streaming-platform/services/upload/pkg/s3"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Server Variables
	port := os.Getenv("PORT")

	//S3 Variables
	region := os.Getenv("REGION")
	accessKeyID := os.Getenv("ACCESKEYID")
	secretAccessKey := os.Getenv("SECRETACCESKEY")
	endpoint := os.Getenv("ENDPOINT")

	//Queue Variables
	endpoint_queue := os.Getenv("ENDPOINT_QUEUE")
	name_queue := os.Getenv("NAME_QUEUE")

	if accessKeyID == "" || secretAccessKey == "" || region == "" ||
		endpoint == "" {
		log.Fatal(
			"missing env: RUSTFS_ACCESS_KEY_ID / " +
				"RUSTFS_SECRET_ACCESS_KEY / RUSTFS_REGION / " +
				"RUSTFS_ENDPOINT_URL",
		)
	}

	client_data := dataStore.Init(accessKeyID, secretAccessKey, region, endpoint)

	queue_client := queue.Init(name_queue, endpoint_queue)

	h := handler.NewHandler(queue_client, client_data)
	mux := http.NewServeMux()
	h.RegisterHandler(mux)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		log.Println("Server is starting...")
		log.Printf("Server läuft auf http://localhost%s\n", srv.Addr)
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
