package handlers

import (
	"encoding/json"
	"io"
	"strings"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// generateShortURL - генерирует короткий URL и сохраняет его в хранилище.
func generateShortURL(originalURL string) (shortURL string) {
	// Сгенерировать короткий id
	shortID := shortener.Shorten(originalURL)
	// Cохранить короткий id в хранилище в RAM
	savedID, added := storage.Set(shortID, originalURL)

	// Если это новый savedID, то есть added == true,
	// то сохранить savedID и оригинальный URL в базу данных и/или в файловое хранилище
	if added {
		app.AddToStore(savedID, originalURL)
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

/*
APIShortenHandler - обрабатывает POST-запросы для создания короткого URL.
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
*/
func APIShortenHandler(w http.ResponseWriter, r *http.Request) {
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

/*
APIShortenBatchHandler - принимает в теле запроса множество URL для сокращения в формате:
```json
[

	{
		"correlation_id": "<строковый идентификатор>",
		"original_url": "<URL для сокращения>"
	},
	...

]
```

В качестве ответа хендлер должен возвращать данные в формате:

```json
[

	{
		"correlation_id": "<строковый идентификатор из объекта запроса>",
		"short_url": "<результирующий сокращённый URL>"
	},
	...

]
```

Все записи о коротких URL сохраняйте в базе данных. Не забудьте добавить реализацию для сохранения в файл и в память.

Стоит помнить, что:

- нужно соблюдать обратную совместимость;
- отправлять пустые батчи не нужно;
- вы умеете сжимать контент по алгоритму gzip;
- изменение в базе можно выполнять в рамках одной транзакции или одного запроса;
- необходимо избегать формирования условий для возникновения состояния гонки (race condition).
*/
func APIShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"` + strings.ReplaceAll(err.Error(), `"`, ` `) + `"}`))
		return
	}

	// Массив для хранения данных запроса
	var req []struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	// Распарсить тело запроса
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Unmarshal error"})
		return
	}

	// Если пустой батч, то вернуть ошибку
	if len(req) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Empty batch"}`))
		return
	}

	// Массив для хранения данных ответа
	resp := make([]struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}, 0, len(req))

	// Обработать каждый элемент массива запроса
	for _, r := range req {
		originalURL := r.OriginalURL
		if originalURL == "" {
			resp = append(resp, struct {
				CorrelationID string `json:"correlation_id"`
				ShortURL      string `json:"short_url"`
			}{CorrelationID: r.CorrelationID, ShortURL: ""})
			continue
		}

		// Сгенерировать короткий id и сохранить его
		shortURL := generateShortURL(originalURL)

		resp = append(resp, struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}{CorrelationID: r.CorrelationID, ShortURL: shortURL})
	}

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
