package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func TestMain(m *testing.M) {
	logger.InitializeLogger()
	config.ParseCommandLine()
	storage.Create()
	filestorage.LoadDataAndLog(config.Params.FileStoragePath)
	db.Connect(3)

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
	storage.PrintKeyValue()
}

func TestAPIShortenHandler(t *testing.T) {
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
	storage.PrintKeyValue()
}

func TestRedirectHandler(t *testing.T) {

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
