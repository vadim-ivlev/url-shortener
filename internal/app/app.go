// Description: Бизнес логика приложения.
// Определенные здесь фунции связывают в себе вызовы функций различных пакетов.
// Цель  -  понизить связанность (coupling) между пакетами.

package app

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/memstore"
	"github.com/vadim-ivlev/url-shortener/internal/pg"
)

// MemStore - хранилище записей URLShortener в оперативной памяти.
var MemStore apptypes.MemStoreInterface

// InitApp инициализирует приложение.
func InitApp() {
	// Инициализировать логгер
	logger.InitializeLogger()

	// Разобрать параметры командной строки
	config.ParseCommandLine()
	// Разобрать переменные окружения
	config.ParseEnv()
	// Вывести параметры конфигурации в лог
	config.PrintParams()

	// Создать хранилище базы данных
	pg.PGStore = pg.New()
	// Подключиться к базе данных
	err := pg.PGStore.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Cannot connect to DB")
		return
	}

	// Создать хранилище в памяти
	MemStore = memstore.New()
	// Почистить хранилища
	//MemStore.Clear()
	filestorage.Clear()

}

// LoadDataToMemStore загружает данные из базы данных в хранилище.
func LoadDataToMemStore() {
	// Загрузить данные из файлового хранилища
	err := LoadFileDataToStorage()
	if err != nil {
		log.Warn().Err(err).Msg("Cannot load data to storage")
	}

	// Загрузить данные из базы данных
	err = LoadDBDataToStorage(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Cannot load data to storage")
	}

	// Печать содержимого хранилища в лог
	MemStore.PrintRecords(5)
}

// Получить короткий URL из shortID
func ShortURL(shortID string) string {
	return config.Params.BaseURL + "/" + shortID
}

// Получить shortID из shortURL
func ShortID(shortURL string) string {
	return strings.TrimPrefix(shortURL, config.Params.BaseURL+"/")
}
