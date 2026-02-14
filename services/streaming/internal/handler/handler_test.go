package handler

import (
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NiklasWillecke/media-platform/services/streaming/internal/cache"
)

func BenchmarkFileServer1(b *testing.B) {
	// Dein Verzeichnis
	dir := http.Dir("./tmp")

	// FileServer-Handler erstellen
	handler := http.FileServer(dir)

	// Beispiel-Request vorbereiten — z. B. eine Datei im tmp-Verzeichnis
	//
	req := make([]*http.Request, 0)
	req1 := httptest.NewRequest(http.MethodGet, "/output000.ts", nil)
	req2 := httptest.NewRequest(http.MethodGet, "/output000.ts", nil)
	req3 := httptest.NewRequest(http.MethodGet, "/output000.ts", nil)
	req4 := httptest.NewRequest(http.MethodGet, "/generated-image.png", nil)
	req5 := httptest.NewRequest(http.MethodGet, "/animal1.jpg", nil)
	req6 := httptest.NewRequest(http.MethodGet, "/animal2.jpg", nil)

	req = append(req, req1, req2, req3, req4, req5, req6)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req[rand.IntN(6)])

		// Antwort auslesen, aber nicht in der Messung bewerten
		io.Copy(io.Discard, rr.Result().Body)
		rr.Result().Body.Close()
	}
}

func BenchmarkFileServer2(b *testing.B) {
	// Dein Verzeichnis

	// FileServer-Handler erstellen
	lru := cache.NewLRUCache(100 * 1024 * 1024)

	// 2. Handler erstellen
	myCache := &MyCache{LRUCache: lru}

	// Beispiel-Request vorbereiten — z. B. eine Datei im tmp-Verzeichnis
	//
	req := make([]*http.Request, 0)
	req1 := httptest.NewRequest(http.MethodGet, "/output000.ts", nil)
	req2 := httptest.NewRequest(http.MethodGet, "/output000.ts", nil)
	req3 := httptest.NewRequest(http.MethodGet, "/output000.ts", nil)
	req4 := httptest.NewRequest(http.MethodGet, "/generated-image.png", nil)
	req5 := httptest.NewRequest(http.MethodGet, "/animal1.jpg", nil)
	req6 := httptest.NewRequest(http.MethodGet, "/animal2.jpg", nil)

	req = append(req, req1, req2, req3, req4, req5, req6)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		myCache.ServeHTTP(rr, req[rand.IntN(6)])

		// Antwort auslesen, aber nicht in der Messung bewerten
		io.Copy(io.Discard, rr.Result().Body)
		rr.Result().Body.Close()
	}
}
