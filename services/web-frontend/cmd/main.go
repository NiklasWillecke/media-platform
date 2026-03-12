package main

import (
	"log"
	"net/http"
	handler "streaming-platform/services/auth/pkg/handler"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Eigener Handler für /login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	})

	http.HandleFunc("/upload", handler.WithAuth(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/upload.html")
	}))

	log.Println("Server läuft auf http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
