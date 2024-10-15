// Description: Функции для загрузки данных из файлового хранилища в storage.

package filestorage

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// LoadDataToStorage - загружает данные из файлового хранилища в storage.
// Возвращает ошибку.
func LoadDataToStorage() (err error) {
	// Открываем файл для чтения
	file, err := os.OpenFile(config.Params.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Warn().Err(err).Msg("Filestorage not found. Probably this is the first launch.")
		return err
	}
	defer file.Close()

	// Читаем все записи из файла
	records := make([]fileStorageRecord, 0)
	decoder := json.NewDecoder(file)
	for {
		var record fileStorageRecord
		if err := decoder.Decode(&record); err != nil {
			break
		}
		records = append(records, record)

		// Извлекаем shortID из record.ShortURL
		shortID := record.ShortURL[len(config.Params.BaseURL)+1:]
		// Добавляем запись в карту хранилища
		storage.Set(shortID, record.OriginalURL)
	}

	log.Info().Msgf("%d Records loaded from filestorage", len(records))
	return nil
}
