package handlers

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/auth"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// generateAndSaveShortURL - генерирует короткий URL и сохраняет его в хранилище.
// Параметры:
// ctx - контекст
// originalURL - оригинальный URL.
// Возвращает:
// shortURL - короткий URL
// aNewOne -  флаг, новый ли это короткий URL. Если true, то это новый короткий URL.
// err - ошибка.
func generateAndSaveShortURL(ctx context.Context, originalURL string) (shortURL string, aNewOne bool, err error) {
	// Сгенерировать короткий id
	shortID := shortener.Shorten(originalURL)
	// Cохранить короткий id в хранилище в RAM
	savedID, aNewOne := storage.Set(shortID, originalURL)

	// Если это новый savedID, то есть aNewOne == true,
	// то сохранить savedID и оригинальный URL в базу данных и/или в файловое хранилище
	if aNewOne {
		err = app.AddToStore(ctx, savedID, originalURL)
	}
	return app.ShortURL(savedID), aNewOne, err
}

// ShortenURLHandler обрабатывает POST-запросы для создания короткого URL.
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	log.Info().Msgf("ShortenURLHandler> User ID '%v' ", userID)

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
	shortURL, aNewOne, err := generateAndSaveShortURL(ctx, app.JoinUserAndURL(userID, originalURL))
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
// При запросе удалённого URL нужно вернуть статус `410 Gone`.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	// если id пустой, то вернуть ошибку
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получить userID из контекста
	userID := GetUserIDFromContext(r.Context())
	log.Info().Msgf("RedirectHandler> User ID from context = '%v' ", userID)

	// Получить оригинальный URL по id и перенаправить
	storedValue := storage.Get(id)
	if storedValue == "" {
		// Проверить не удаленный ли это URL
		if storage.IsDeletedKey(id) {
			http.Error(w, "URL was deleted", http.StatusGone)
			return
		}
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	log.Info().Msgf("RedirectHandler> storedValue = '%v'", storedValue)
	storedUserID, storedURL := app.SplitUserAndURL(storedValue)
	log.Info().Msgf("RedirectHandler> storedUserID = '%v', storedURL = '%v'", storedUserID, storedURL)

	http.Redirect(w, r, storedURL, http.StatusTemporaryRedirect)
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
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	log.Info().Msgf("APIShortenHandler> User ID '%v' ", userID)

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
	shortURL, aNewOne, err := generateAndSaveShortURL(ctx, app.JoinUserAndURL(userID, originalURL))
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

	// Установить статус ответа в зависимости от наличия записи
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
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	log.Info().Msgf("APIShortenBatchHandler> User ID '%v' ", userID)

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
		shortURL, _, err := generateAndSaveShortURL(ctx, app.JoinUserAndURL(userID, originalURL))
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

// dataRec - Тип записи выходных данных для APIUserURLsHandler
type dataRec struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

/*
APIUserURLsHandler - возвращает пользователю все когда-либо сокращённые им `URL` в формате:
```json
[

	{
		"short_url": "http://...",
		"original_url": "http://..."
	},
	...

]
```

- Если кука не содержит `ID` пользователя, хендлер должен возвращать HTTP-статус `401 Unauthorized`.
- При отсутствии сокращённых пользователем URL хендлер должен отдавать HTTP-статус `204 No Content`.
*/
func APIUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	newUserID := GetNewUserIDFromContext(r.Context())

	// Проверить, что ID пользователя не пустой, или это новый пользователь с только что сгенерированным ID
	if userID == "" || newUserID == "new" {
		log.Error().Msg("APIUserURLsHandler> User ID not found or User ID was generated on the fly")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Unauthorized: No user ID"}`))
		return
	}

	// Получить все короткие URL пользователя
	urls := app.GetUserURLs(userID)

	// Подготовливаем тело ответа
	respBody, err := json.Marshal(urls)
	if err != nil {
		log.Error().Err(err).Msg("APIUserURLsHandler> Marshal error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Marshal error"}`))
		return
	}
	log.Info().Msg("---------------------------")
	log.Info().Msgf("APIUserURLsHandler> Response: %v", string(respBody))
	log.Info().Msg("---------------------------")

	// Устанавливаем статус ответа в зависимости от наличия записей
	status := http.StatusOK
	// Если коротких URL нет, то вернуть статус 204
	if len(urls) == 0 {
		log.Warn().Msg("APIUserURLsHandler> No content")
		status = http.StatusNoContent
		// status = http.StatusUnauthorized
	}

	// Отправляем ответ
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

// GetUserIDFromContext - получает ID пользователя из контекста.
func GetUserIDFromContext(ctx context.Context) (userID string) {
	userID, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		log.Error().Msg("GetUserIDFromContext> User ID not found in context")
	}
	// log.Info().Msgf("GetUserIDFromContext> User ID '%v' ", userID)
	return userID
}

// GetNewUserIDFromContext - получает новый ли это (сгенерированный налету) ID пользователя из контекста.
// Если это новый ID, то возвращает "new".
func GetNewUserIDFromContext(ctx context.Context) (newUserID string) {
	newUserID, ok := ctx.Value(auth.NewUserIDKey).(string)
	if !ok {
		log.Error().Msg("GetNewUserIDFromContext> New User ID not found in context")
	}
	log.Info().Msgf("GetNewUserIDFromContext> New User ID flag '%v' ", newUserID)
	return newUserID
}

// APIDeleteURLsHandler - обрабатывает DELETE-запросы для удаления коротких URL.
// в теле запроса принимает список идентификаторов сокращённых URL для асинхронного удаления.
// Запрос может быть таким:
// ```http
// DELETE http://localhost:8080/api/user/urls
// Content-Type: application/json
//
// ["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]
// ```
//
// В случае успешного приёма запроса хендлер должен возвращать HTTP-статус `202 Accepted`.
// Фактический результат удаления может происходить позже.
// Оповещать пользователя об успешности или неуспешности не нужно.
func APIDeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	// Получить userID из контекста
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	log.Info().Msgf("APIDeleteURLsHandler> User ID '%v' ", userID)

	// Прочитать тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"` + strings.ReplaceAll(err.Error(), `"`, ` `) + `"}`))
		return
	}

	// Массив идентификаторов для удаления
	ids := []any{}

	// Распарсить тело запроса в массив идентификаторов
	err = json.Unmarshal(body, &ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		errorText := strings.ReplaceAll(err.Error(), `"`, ` `)
		json.NewEncoder(w).Encode(map[string]string{"error": errorText})
		return
	}

	// Если массив идентификаторов пустой, то вернуть ошибку
	if len(ids) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Empty batch"}`))
		return
	}

	// Удалить короткие URL
	go deleteKeys(ctx, userID, ids)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Accepted"}`))
}

// deleteKeys - удаляет короткие URL из хранилища в RAM и из базы данных.
func deleteKeys(ctx context.Context, userID string, ids []any) error {
	err := storage.DeleteKeys(userID, ids)
	if err != nil {
		log.Warn().Err(err).Msg("APIDeleteURLsHandler> Cannot delete shortIDs from RAM")
		return err
	}

	err = app.DeleteKeysFromStore(ctx, userID, ids)
	if err != nil {
		log.Warn().Err(err).Msg("APIDeleteURLsHandler> Cannot delete shortIDs from the database")
	}

	return nil
}
