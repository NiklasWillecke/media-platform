package main

import (
	"log"
	"net/http"
	handler "streaming-platform/services/auth/pkg/handler"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Eigener Handler für /login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	})

	//With Auth middleware
	http.HandleFunc("/upload", handler.WithAuth(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, "./static/upload.html")
	}))

	log.Println("Server läuft auf http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
