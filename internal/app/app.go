// Бизнес логика приложения.
// Определенные здесь фунции связывают в себе вызовы функций различных пакетов.
// Цель  -  понизить связанность (coupling) между пакетами.

package app

import (
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

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
		loadDataFromDB()
	case config.Params.FileStoragePath != "":
		filestorage.LoadDataAndLog(config.Params.FileStoragePath)
	default:
		log.Info().Msg("LoadData(). No persistent data store specified")
	}
}

// loadDataFromDB - загружает данные из базы данных в storage.
func loadDataFromDB() {
	if !db.IsConnected() {
		log.Error().Msg("loadDataFromDB(). No connection to DB")
		return
	}
	data, err := db.GetData()
	if err != nil {
		log.Warn().Err(err).Msg("loadDataFromDB(). Cannot get data from DB")
	}
	storage.LoadData(data)
}
