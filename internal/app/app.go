// Description: Бизнес логика приложения.
// Определенные здесь фунции связывают в себе вызовы функций различных пакетов.
// Цель  -  понизить связанность (coupling) между пакетами.

package app

import (
	"context"
	"strings"

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
	// Разобрать переменные окружения
	config.ParseEnv()
	// Вывести параметры конфигурации в лог
	config.PrintParams()

	// Создать хранилище в памяти
	storage.Create()

	// Подключиться к базе данных с 1-й попытки
	db.TryToConnect(1)
	// Выполнить миграции базы данных
	db.MigrateUp("./migrations")

	// Загрузить данные из базы данных или из файлового хранилища в storage
	err := LoadDataToStorage(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg("Cannot load data to storage")
	}
	// Печать содержимого хранилища в лог
	storage.PrintContent(0)
}

// Получить короткий URL из shortID
func ShortURL(shortID string) string {
	return config.Params.BaseURL + "/" + shortID
}

// Получить shortID из shortURL
func ShortID(shortURL string) string {
	return strings.TrimPrefix(shortURL, config.Params.BaseURL+"/")
}

// LoadDataToStorage - загружает данные из базы данных или из файлового хранилища в storage.
// Если указана DatabaseDSN в конфигурации, то загрузить данные из базы данных.
// В противном случае, если указан FileStoragePath в конфигурации, то загрузить данные из файлового хранилища.
// Если ни один из параметров не указан, то ничего не загружать.
// Параметры:
// - ctx - контекст
// Возвращает ошибку, если загрузка данных не удалась.
func LoadDataToStorage(ctx context.Context) (err error) {
	switch {
	case config.Params.DatabaseDSN != "":
		err = LoadDBDataToStorage(ctx)
	case config.Params.FileStoragePath != "":
		err = LoadFileDataToStorage()
	default:
		log.Info().Msg("LoadData(). No persistent data store specified")
	}
	return err
}

// AddToStore сохраняет короткий и оригинальный URL в базу данных или в файловое хранилище.
// Если указана DatabaseDSN в конфигурации, то сохранять данные в базу данных.
// В противном случае, если указан FileStoragePath в конфигурации, то сохранять данные в файловое хранилище.
// Если ни один из параметров не указан, то ничего не сохранять.
// Параметры:
// - ctx - контекст
// - shortID - короткий ID
// - originalURL - оригинальный URL
// Возвращает ошибку, если сохранение не удалось.
func AddToStore(ctx context.Context, shortID, originalURL string) (err error) {
	switch {
	case config.Params.DatabaseDSN != "":
		// сохранить shortID и оригинальный URL в базу данных
		err = db.Store(ctx, shortID, originalURL)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortID in the database")
			return err
		}
	case config.Params.FileStoragePath != "":
		// сохранить shortID и оригинальный URL в файловое хранилище
		err := filestorage.Store(ShortURL(shortID), originalURL)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortened url in the filestorage")
			return err
		}
	default:
		log.Info().Msg("AddToStore(). No persistent data store specified")
	}
	return nil
}

// GetUserURLs возвращает пользователю все когда-либо сокращённые им `URL` в формате:
// ```json
// [
//
//	{
//	    "short_url": "http://...",
//	    "original_url": "http://..."
//	},
//	...
//
// ]
func GetUserURLs(userID string) (urls []map[string]string) {
	urls = make([]map[string]string, 0)
	// Получить данные из хранилища
	storageData := storage.GetData()
	for recordShortID, recordValue := range storageData {

		recordUserID, RecordOriginalURL := SplitUserAndURL(recordValue)
		if recordUserID != userID {
			continue
		}
		urls = append(urls, map[string]string{"short_url": ShortURL(recordShortID), "original_url": RecordOriginalURL})
	}
	return urls
}

// JoinUserAndURL - объединяет ID пользователя и URL.
func JoinUserAndURL(userID, URL string) string {
	// return URL
	return userID + "@" + URL
}

// SplitUserAndURL - разделяет ID пользователя и URL.
func SplitUserAndURL(userAndURL string) (userID, URL string) {
	parts := strings.Split(userAndURL, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func DeleteKeysFromStore(ctx context.Context, userID string, keys []any) error {
	// Пометить keys как удаленные в базе данных
	if config.Params.DatabaseDSN != "" {
		err := db.DeleteKeys(ctx, userID, keys)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot delete shortID from the database")
			return err
		}
	}
	// Пометить keys как удаленные в файловом хранилище
	if config.Params.FileStoragePath != "" {
		err := DumpDataToFilestorage()
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortened url in the filestorage")
			return err
		}
	}
	return nil
}

// DumpDataToFilestorage - сохраняет данные из RAM в файловое хранилище.
func DumpDataToFilestorage() error {
	if config.Params.FileStoragePath == "" {
		return nil
	}
	filestorage.Clear()
	storageData := storage.GetData()
	for recordShortID, recordValue := range storageData {
		err := filestorage.Store(ShortURL(recordShortID), recordValue)
		if err != nil {
			log.Warn().Err(err).Msg("Cannot save shortened url in the filestorage")
			return err
		}
	}
	return nil
}
