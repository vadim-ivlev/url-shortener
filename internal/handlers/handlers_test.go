package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		log.Info().Msg("Skipping testing in CI environment")
		t.Skip("Skipping testing in CI environment")
	}
}

func TestMain(m *testing.M) {
	// Перейти в корневую директорию проекта
	os.Chdir("../../")

	app.InitApp()

	InitTestTable()
	os.Exit(m.Run())
}

type want struct {
	postReturnCode int
	getReturnCode  int
	shortURL       string
	contentType    string
}

type testTable = []struct {
	name string
	url  string
	want want
}

var tests testTable

func InitTestTable() {
	tests = testTable{
		{
			name: "Empty",
			url:  "",
			want: want{
				postReturnCode: http.StatusBadRequest,
				getReturnCode:  http.StatusBadRequest,
				shortURL:       "Empty URL",
				contentType:    "text/plain",
			},
		},
		{
			name: "Google",
			url:  "https://www.google.com",
			want: want{
				postReturnCode: http.StatusCreated,
				getReturnCode:  http.StatusTemporaryRedirect,
				shortURL:       config.Params.BaseURL + "/F870F1E9",
				contentType:    "text/plain",
			},
		},
		{
			name: "Youtube",
			url:  "https://www.youtube.com",
			want: want{
				postReturnCode: http.StatusCreated,
				getReturnCode:  http.StatusTemporaryRedirect,
				shortURL:       config.Params.BaseURL + "/4AED1C05",
				contentType:    "text/plain",
			},
		},
		{
			name: "Google2",
			url:  "https://www.google.com",
			want: want{
				postReturnCode: http.StatusCreated,
				getReturnCode:  http.StatusTemporaryRedirect,
				shortURL:       config.Params.BaseURL + "/F870F1E9",
				contentType:    "text/plain",
			},
		},
	}
}

func TestShortenURLHandler(t *testing.T) {
	skipCI(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			rec := httptest.NewRecorder()
			// Вызов обработчика
			ShortenURLHandler(rec, req)
			// Проверка ответа
			assert.Equal(t, tt.want.postReturnCode, rec.Code)
			bodyString := strings.TrimSpace(rec.Body.String())
			assert.Equal(t, tt.want.shortURL, bodyString)
			assert.Contains(t, rec.Header().Get("Content-Type"), tt.want.contentType)

		})
	}
	storage.PrintContent(3)
}

func TestAPIShortenHandler(t *testing.T) {
	skipCI(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"`+tt.url+`"}`))
			rec := httptest.NewRecorder()
			// Вызов обработчика
			APIShortenHandler(rec, req)
			// Проверка ответа
			assert.Equal(t, tt.want.postReturnCode, rec.Code)
			bodyString := strings.TrimSpace(rec.Body.String())
			assert.Contains(t, bodyString, tt.want.shortURL)
			assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
			fmt.Printf("Content-Type: %v\n", rec.Header().Get("Content-Type"))
		})
	}
	storage.PrintContent(3)
}

func TestRedirectHandler(t *testing.T) {
	skipCI(t)
	// Добавим тесты для проверки перенаправления
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			rec := httptest.NewRecorder()
			// Вызов обработчика
			ShortenURLHandler(rec, req)
		})
	}

	// Проверим перенаправление
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := getID(tt.want.shortURL)
			req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
			rec := httptest.NewRecorder()

			// Добавим параметр в URL используя контекст
			req = WithURLParam(req, "id", id)

			// Вызов обработчика
			RedirectHandler(rec, req)

			// Проверка ответа
			assert.Equal(t, tt.want.getReturnCode, rec.Code)
			assert.Equal(t, tt.url, rec.Header().Get("Location"))
			fmt.Printf("id = %v, Location = %v, Code = %v\n", id, rec.Header().Get("Location"), rec.Code)
		})
	}
}

// WithURLParam возвращает указатель на объект запроса
// с добавленными URL-параметрами в новом объекте chi.Context.
// https://haykot.dev/blog/til-testing-parametrized-urls-with-chi-router/
// https://github.com/go-chi/chi/issues/76
func WithURLParam(r *http.Request, key, value string) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add(key, value)
	newCtx := context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx)
	req := r.WithContext(newCtx)
	return req
}

// возвращает послледнюю часть URL (после 22 символа) или "" если URL короче
func getID(url string) (id string) {
	if len(url) > 22 {
		id = url[22:]
	}
	fmt.Printf("getId()   : '%v'\n", id)
	return
}

func TestPingHandler(t *testing.T) {
	skipCI(t)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	// подключенная БД
	db.Connect(3)
	rec := httptest.NewRecorder()
	PingHandler(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// отключенная БД
	db.Disconnect()
	rec = httptest.NewRecorder()
	PingHandler(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

/*
APIShortenBatchHandler принимает в теле запроса множество URL для сокращения в формате:
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
func aTestAPIShortenBatchHandler(t *testing.T) {
	skipCI(t)

	emptyBodyString := ``
	noElementsBodyString := ` [] `
	normalBodyString := `
		[
			{
				"correlation_id": "0",
				"original_url": ""
			},
			{
				"correlation_id": "1",
				"original_url": "https://www.google.com"
			},
			{
				"correlation_id": "2",
				"original_url": "https://www.youtube.com"
			}
		]	
	`

	type args struct {
		bodyString string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Empty",
			args: args{
				bodyString: emptyBodyString,
			},
		},
		{
			name: "NoElements",
			args: args{
				bodyString: noElementsBodyString,
			},
		},
		{
			name: "Normal",
			args: args{
				bodyString: normalBodyString,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.args.bodyString))
			rec := httptest.NewRecorder()
			APIShortenBatchHandler(rec, req)

			// Проверка статуса ответа
			status := rec.Code
			log.Info().Msgf("Status: %v", status)
			assert.Equal(t, http.StatusOK, status)

			// Проверка типа контента
			contentType := rec.Header().Get("Content-Type")
			log.Info().Msgf("Content-Type: %v", contentType)
			assert.Equal(t, "application/json", contentType)

			// Проверка тела ответа
			log.Info().Msgf("Body: %v", rec.Body.String())

			// Распарсить тело ответа как в массив структур
			resp := []struct {
				CorrelationID string `json:"correlation_id"`
				ShortURL      string `json:"short_url"`
			}{}
			err := json.Unmarshal([]byte(rec.Body.String()), &resp)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			log.Info().Msgf("Response: %+v", resp)

			// Проверка количества элементов в ответе
			assert.Equal(t, 3, len(resp))
			// Проверка корреляционных идентификаторов
			assert.Equal(t, "0", resp[0].CorrelationID)
			assert.Equal(t, "1", resp[1].CorrelationID)
			assert.Equal(t, "2", resp[2].CorrelationID)
		})
	}
}
