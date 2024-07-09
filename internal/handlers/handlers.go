package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// ShortenURLHandler обрабатывает POST-запросы для создания короткого URL.
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	// Сгенерировать короткий id и сохранить его
	shortID := shortener.Shorten(originalURL)
	savedID := storage.Set(shortID, originalURL)
	// Сгенерировать короткий URL
	shortURL := "http://localhost:8080/" + savedID

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortURL))
}

// RedirectHandler обрабатывает GET-запросы для перенаправления на оригинальный URL.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// chiShortID := chi.URLParam(r, "id")
	// shortID := r.URL.Path[1:]
	// fmt.Printf("   shortID: '%s'\n", shortID)
	// fmt.Printf("chiShortID: '%s'\n", chiShortID)

	shortID := chi.URLParam(r, "id")

	if shortID == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получить оригинальный URL и перенаправить
	originalURL := storage.Get(shortID)
	if originalURL == "" {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
