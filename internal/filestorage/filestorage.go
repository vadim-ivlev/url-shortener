// Description: Файловое хранилище для хранения записей в формате JSON.

package filestorage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

// createDirIfNotExists - создает директорию в которой будет храниться файл хранилища, если ее нет.
// Параметры:
// - filePath - путь к файлу хранилища.
func createDirIfNotExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// Clear - очищает файловое хранилище.
func Clear() error {
	if !config.UseFileStorage() {
		return nil
	}
	return os.Remove(config.Params.FileStoragePath)
}

// AddRecord - добавляет запись в файловое хранилище.
//
// Параметры:
// - record - запись для сохранения.
//
// Возвращает ошибку, если запись не удалась.
func AddRecord(record apptypes.URLShortener) error {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseFileStorage() {
		return nil
	}

	// Преобразуем запись в JSON
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// Создаем директорию для файла хранилища, если ее нет
	if err := createDirIfNotExists(config.Params.FileStoragePath); err != nil {
		return err
	}

	// Открываем файл для записи (добавляем в конец файла) или создаем новый
	file, err := os.OpenFile(config.Params.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем recordJSON  в файл
	if _, err := file.Write(append(recordJSON, '\n')); err != nil {
		return err
	}
	log.Info().Msgf("Record saved to filestorage: %s in file %s", recordJSON, config.Params.FileStoragePath)
	return nil
}

// DumpRecords - сохраняет записи в файловое хранилище.
//
// Параметры:
// - records - записи для сохранения.
//
// Возвращает ошибку, если запись не удалась.
func DumpRecords(records []apptypes.URLShortener) error {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseFileStorage() {
		return nil
	}

	// Создаем директорию для файла хранилища, если ее нет
	if err := createDirIfNotExists(config.Params.FileStoragePath); err != nil {
		return err
	}

	// Открываем файл для записи или создаем новый
	// file, err := os.OpenFile(config.Params.FileStoragePath, os.O_CREATE|os.O_WRONLY| os.O_TRUNC , 0644)
	file, err := os.Create(config.Params.FileStoragePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, record := range records {
		// Преобразуем запись в JSON
		recordJSON, err := json.Marshal(record)
		if err != nil {
			return err
		}

		// Записываем recordJSON  в файл
		if _, err := file.Write(append(recordJSON, '\n')); err != nil {
			return err
		}
	}
	return nil
}
