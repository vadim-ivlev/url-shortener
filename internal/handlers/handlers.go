package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadim-ivlev/url-shortener/internal/config"
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
	shortURL := config.BaseURL + "/" + savedID

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortURL))
}

/*
APIShortenHandler обрабатывает POST-запросы для создания короткого URL.
Обслуживает эндпоинт POST /api/shorten,
принимает в теле запроса JSON-объект `{"url":"<some_url>"}`
и возвращает в ответ объект `{"result":"<short_url>"}`.
Запрос может иметь такой вид:

	POST http://localhost:8080/api/shorten HTTP/1.1
	Host: localhost:8080
	Content-Type: application/json

	{
	"url": "https://practicum.yandex.ru"
	}

Ответ может быть таким:

	HTTP/1.1 201 OK
	Content-Type: application/json
	Content-Length: 30

	{
	"result": "http://localhost:8080/EwHXdJfB"
	}

При реализации задействуйте пакеты:

	encoding/json
*/
func APIShortenHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	type apiShortenRequest struct {
		URL string `json:"url"`
	}

	var req apiShortenRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unmarshal error"})
		return
	}

	originalURL := req.URL
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Empty URL"}`))
		return
	}

	// Сгенерировать короткий id и сохранить его
	shortID := shortener.Shorten(originalURL)
	savedID := storage.Set(shortID, originalURL)
	// Сгенерировать короткий URL
	shortURL := config.BaseURL + "/" + savedID

	type apiShortenResponse struct {
		Result string `json:"result"`
	}

	resp := apiShortenResponse{Result: shortURL}
	respBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

// RedirectHandler обрабатывает GET-запросы для перенаправления на оригинальный URL.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	// если id пустой, то вернуть ошибку
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получить оригинальный URL по id и перенаправить
	originalURL := storage.Get(id)
	if originalURL == "" {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
