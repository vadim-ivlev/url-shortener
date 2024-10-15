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

// generateAndSaveShortURL - генерирует короткий URL и сохраняет его в хранилище.
// originalURL - оригинальный URL.
// Возвращает:
// shortURL - короткий URL
// aNewOne -  флаг, новый ли это короткий URL. Если true, то это новый короткий URL.
// err - ошибка.
func generateAndSaveShortURL(originalURL string) (shortURL string, aNewOne bool, err error) {
	// Сгенерировать короткий id
	shortID := shortener.Shorten(originalURL)
	// Cохранить короткий id в хранилище в RAM
	savedID, aNewOne := storage.Set(shortID, originalURL)

	// Если это новый savedID, то есть aNewOne == true,
	// то сохранить savedID и оригинальный URL в базу данных и/или в файловое хранилище
	if aNewOne {
		err = app.AddToStore(savedID, originalURL)
	}
	return app.ShortURL(savedID), aNewOne, err
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
	shortURL, aNewOne, err := generateAndSaveShortURL(originalURL)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Определить статус ответа
	status := http.StatusCreated
	// Если короткий URL уже существует, то вернуть статус 409
	if !aNewOne {
		status = http.StatusConflict
	}

	w.WriteHeader(status)
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
		http.Error(w, "No connection do database", http.StatusInternalServerError)
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
	shortURL, aNewOne, err := generateAndSaveShortURL(originalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"` + strings.ReplaceAll(err.Error(), `"`, ` `) + `"}`))
		return
	}

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

	// Определить статус ответа
	status := http.StatusCreated
	// Если короткий URL уже существует, то вернуть статус 409
	if !aNewOne {
		status = http.StatusConflict
	}

	// Отправляем ответ
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

// Типы входных и выходных данных для APIShortenBatchHandler **********************

// Тип записи входных данных
type inpRec struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Тип записи выходных данных
type outRec struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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
	// Прочитать тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"` + strings.ReplaceAll(err.Error(), `"`, ` `) + `"}`))
		return
	}

	// Массив входных данных запроса
	inputRecords := []inpRec{}

	// Распарсить тело запроса в массив входных данных
	err = json.Unmarshal(body, &inputRecords)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		errorText := strings.ReplaceAll(err.Error(), `"`, ` `)
		json.NewEncoder(w).Encode(map[string]string{"error": errorText})
		return
	}

	// Если массив входных данных пустой, то вернуть ошибку
	if len(inputRecords) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Empty batch"}`))
		return
	}

	// Массив выходных данных ответа c емкостью равной длине входного массива
	outputRecords := make([]outRec, 0, len(inputRecords))

	// Обработать каждый элемент массива входных данных
	for _, r := range inputRecords {
		originalURL := r.OriginalURL

		// Если originalURL пустой, то вернуть пустой shortURL, не сохраняя его в хранилище и БД
		if originalURL == "" {
			outputRecords = append(outputRecords, outRec{CorrelationID: r.CorrelationID, ShortURL: ""})
			continue
		}
		// Сгенерировать короткий id и сохранить его в хранилище и в БД
		shortURL, _, err := generateAndSaveShortURL(originalURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error":"` + strings.ReplaceAll(err.Error(), `"`, ` `) + `"}`))
			return
		}

		outputRecords = append(outputRecords, outRec{CorrelationID: r.CorrelationID, ShortURL: shortURL})
	}

	// Подготовливаем тело ответа
	respBody, err := json.Marshal(outputRecords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Marshal error"}`))
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}
