// Description: Файловое хранилище для хранения записей в формате JSON.
// Пример содержимого файла хранилища:
// ```json
// {"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}
// {"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
// {"uuid":"3","short_url":"dG56Hqxm","original_url":"http://practicum.yandex.ru"}
// ```

package filestorage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

// FileStorageRecord - структура для хранения записи в файловом хранилище.
type FileStorageRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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

// Store - сохраняет данные в файловое хранилище.
// Параметры:
// - shortURL - укороченный URL.
// - originalURL - оригинальный URL.
// Возвращает ошибку, если запись не удалась.
func Store(shortURL, originalURL string) error {
	// Генерируем новый UUID
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	// Создаем новую запись
	record := FileStorageRecord{
		UUID:        uuid.String(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
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
	return nil
}
