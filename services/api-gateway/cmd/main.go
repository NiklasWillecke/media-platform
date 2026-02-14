package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	db "github.com/niklaswillecke/streaming-platform/shared/db/generated"
)

// 1. Checks if Video url is valid
// 2. Checks if user is logged in and have permissions
// 3. Find best Server to route user into
// 4.
type App struct {
	DB *pgx.Conn
}

func (h *App) handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path[1:])
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

	queries := db.New(h.DB)
	queries.GetUnlock(context.Background(), 1)
}

func main() {

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://postgres:example@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close(ctx)

	newDB := &App{DB: conn}

	mux := http.NewServeMux()

	mux.HandleFunc("/", newDB.handler)
	log.Fatal(http.ListenAndServe(":8080", mux))

	path := "/1"
	userID := "user123"
	signingKey := "dev-secret-key"
	expires := time.Now().Add(1 * time.Hour).Unix()

	data := fmt.Sprintf("%s:%s:%d", path, userID, expires)
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write([]byte(data))
	token := hex.EncodeToString(mac.Sum(nil))

	url := fmt.Sprintf("http://localhost:8080/hls%s?token=%s&expires=%d&uid=%s",
		path, token, expires, userID)

	fmt.Println("Test URL:")
	fmt.Println(url)
}
