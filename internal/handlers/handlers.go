package handlers

import (
	"encoding/json"
	"io"
	"strings"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// AddToStore сохраняет короткий и оригинальный URL в базу данных и/или в файловое хранилище.
// Если указана DatabaseDSN в конфигурации, то сохранять данные в базу данных.
// В противном случае, если указан FileStoragePath в конфигурации, то сохранять данные в файловое хранилище.
// Если ни один из параметров не указан, то ничего не сохранять.
func AddToStore(shortID, originalURL string) {
	switch {
	case config.Params.DatabaseDSN != "":
		// сохранить shortID и оригинальный URL в базу данных
		err := db.Store(shortID, originalURL)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortID in the database")
		}
	case config.Params.FileStoragePath != "":
		// сохранить shortID и оригинальный URL в файловое хранилище
		err := filestorage.Store(config.Params.FileStoragePath, app.ShortURL(shortID), originalURL)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortened url in the filestorage")
		}
	default:
		log.Info().Msg("AddToStore(). No persistent data store specified")
	}
}

// generateShortURL - генерирует короткий URL и сохраняет его в хранилище.
func generateShortURL(originalURL string) (shortURL string) {
	// Сгенерировать короткий id
	shortID := shortener.Shorten(originalURL)
	// и сохранить его в хранилище в RAM
	savedID, added := storage.Set(shortID, originalURL)

	// Если это новый savedID, то есть added == true,
	// то сохранить savedID и оригинальный URL в базу данных и/или в файловое хранилище
	if added {
		AddToStore(savedID, originalURL)
	}
	return app.ShortURL(shortID)
}

// ShortenURLHandler обрабатывает POST-запросы для создания короткого URL.
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
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
	shortURL := generateShortURL(originalURL)

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
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"` + strings.ReplaceAll(err.Error(), `"`, ` `) + `"}`))
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Unmarshal error"})
		return
	}

	originalURL := req.URL
	if originalURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Empty URL"}`))
		return
	}

	// Сгенерировать короткий id и сохранить его
	shortURL := generateShortURL(originalURL)

	resp := struct {
		Result string `json:"result"`
	}{Result: shortURL}

	respBody, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Marshal error"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
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

// PingHandler - при запросе проверяет соединение с базой данных.
// При успешной проверке хендлер должен вернуть HTTP-статус `200 OK`, при неуспешной — `500 Internal Server Error`.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	if db.IsConnected() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
