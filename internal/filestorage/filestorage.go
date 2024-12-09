// Description: Файловое хранилище для хранения записей в формате JSON.
// Пример содержимого файла хранилища:
// ```json
// {"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}
// {"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
// {"uuid":"3","short_url":"dG56Hqxm","original_url":"http://practicum.yandex.ru"}
// ```

package filestorage

import (
	"os"
	"path/filepath"

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
