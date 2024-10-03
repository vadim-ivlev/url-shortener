// Description: Бизнес логика приложения.
// Определенные здесь фунции связывают в себе вызовы функций различных пакетов.
// Цель  -  понизить связанность (coupling) между пакетами.

package app

import (
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// InitApp инициализирует приложение.
func InitApp() {
	// Инициализировать логгер
	logger.InitializeLogger()
	// Разобрать параметры командной строки
	config.ParseCommandLine()
	// Вывести параметры конфигурации в лог
	config.PrintParams()

	// Создать хранилище в памяти
	storage.Create()

	// Подключиться к базе данных
	db.Connect(1)
	// Выполнить миграции базы данных
	db.MigrateUp("./migrations")

	// Загрузить данные из базы данных или из файлового хранилища в storage
	LoadDataToStorage()
	// Печать содержимого хранилища в лог
	storage.PrintContent(0)
}

// Получить короткий URL из shortID
func ShortURL(shortID string) string {
	return config.Params.BaseURL + "/" + shortID
}

// Получить shortID из shortURL
func ShortID(shortURL string) string {
	return shortURL[len(config.Params.BaseURL)+1:]
}

// LoadDataToStorage - загружает данные из базы данных и/или из файлового хранилища в storage.
// Если указана DatabaseDSN в конфигурации, то загрузить данные из базы данных.
// В противном случае, если указан FileStoragePath в конфигурации, то загрузить данные из файлового хранилища.
// Если ни один из параметров не указан, то ничего не загружать.
func LoadDataToStorage() {
	switch {
	case config.Params.DatabaseDSN != "":
		db.LoadData()
	case config.Params.FileStoragePath != "":
		filestorage.LoadData()
	default:
		log.Info().Msg("LoadData(). No persistent data store specified")
	}
}

// AddToStore сохраняет короткий и оригинальный URL в базу данных или в файловое хранилище.
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
		err := filestorage.Store(ShortURL(shortID), originalURL)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortened url in the filestorage")
		}
	default:
		log.Info().Msg("AddToStore(). No persistent data store specified")
	}
}
