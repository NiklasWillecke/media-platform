package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/NiklasWillecke/media-platform/services/streaming/internal/cache"
)

type MyCache struct {
	*cache.LRUCache
}

type HLSServer struct {
	contentDir string
	signingKey []byte
}

func StartFileServer() {
	directoryPath := "./tmp"

	// Check if the directory exists
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found.\n", directoryPath)
		return
	}

	// Create a file server handler to serve the directory's contents
	fileServer := http.FileServer(http.Dir(directoryPath))

	// Create a new HTTP server and handle requests
	http.Handle("/", fileServer)

	// Start the server on port 8080
	port := 8080
	fmt.Printf("Server started at http://localhost:%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func (m *MyCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("./tmp", r.URL.Path)

	data, ok := m.Get(path)

	if ok {
		// Datei im Cache -> direkt ausliefern
		http.ServeContent(w, r, path, time.Now(), bytes.NewReader(data))
		return
	}

	// sonst Datei laden
	fileData, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// in Cache speichern
	m.Set(path, fileData)

	http.ServeContent(w, r, path, time.Now(), bytes.NewReader(fileData))
}
