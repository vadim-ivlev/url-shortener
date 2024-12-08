package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

// DB - пул соединений с базой данных
var DB *sqlx.DB = nil

// CreatePool - создает пул соединений с базой данных
func CreatePool() (err error) {
	DB, err = sqlx.Connect("postgres", config.Params.DatabaseDSN)
	return err
}

// TryToConnect - Пытается соединиться с базой данных повторяя попытки в случае неудачи.
// numAttempts - количество попыток
func TryToConnect(numAttempts int) (err error) {
	err = errors.New("no attempts to connect to DB")
	for i := 1; i <= numAttempts; i++ {
		err = CreatePool()
		if err == nil {
			log.Info().Msg("Connected to DB")
			return err
		}
		log.Warn().Err(err).Msgf("Waiting for db connection. Attempt # %d", i)
		time.Sleep(time.Second)
	}
	log.Error().Msg("Failed to connect to DB")
	return err
}

// Disconnect - закрывает соединение с базой данных
func Disconnect() {
	if DB != nil {
		DB.Close()
	}
	DB = nil
}

// IsConnected - проверяет, установлено ли соединение с базой данных
func IsConnected() bool {
	return DB != nil && DB.Ping() == nil
}

// Store - сохраняет данные в базу данных.
// Параметры:
// - ctx - контекст
// - shortID - укороченный ID.
// - originalURL - оригинальный URL.
// Возвращает ошибку, если запись не удалась.
func Store(ctx context.Context, shortID, originalURL string) error {
	if !IsConnected() {
		return errors.New("Store. No connection to DB")
	}
	_, err := DB.ExecContext(ctx, "INSERT INTO urls (short_id, original_url) VALUES ($1, $2)", shortID, originalURL)
	return err
}

// AddRecord - добавляет запись в базу данных.
//
// Параметры:
// - record - запись для сохранения.
//
// Возвращает ошибку, если запись не удалась.
func AddRecord(record apptypes.URLShortener) error {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseDatabase() {
		return nil
	}

	if !IsConnected() {
		return errors.New("AddRecord. No connection to DB")
	}

	_, err := DB.Exec("INSERT INTO urls1 (idx, short_id, original_url, user_id, deleted) VALUES ($1, $2, $3, $4, $5)", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted)
	return err
}

// UpdateRecord - обновляет запись в базе данных.
//
// Параметры:
// - record - запись для сохранения.
//
// Возвращает ошибку, если запись не удалась.
func UpdateRecord(record apptypes.URLShortener) error {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseDatabase() {
		return nil
	}

	if !IsConnected() {
		return errors.New("UpdateRecord. No connection to DB")
	}

	_, err := DB.Exec("UPDATE urls1 SET idx = $1,  short_id = $2, original_url = $3, user_id = $4, deleted = $5 WHERE short_id = $6", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted, record.ShortID)
	return err
}

// Clear - очищает таблицу urls
// - ctx - контекст
// Возвращает ошибку, если очистка не удалась.
func Clear(ctx context.Context) error {
	if !IsConnected() {
		return errors.New("Clear. No connection to DB")
	}
	_, err := DB.ExecContext(ctx, "DELETE FROM urls")
	return err
}

// GetData - возвращает данные из базы данных в виде map[string]string,
// где ключ - short_id, значение - original_url.
// Параметры:
// - ctx - контекст
func GetData(ctx context.Context) (data map[string]string, err error) {
	if !IsConnected() {
		return nil, errors.New("GetData. No connection to DB")
	}

	rows, err := DB.QueryxContext(ctx, "SELECT short_id, original_url FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	data = make(map[string]string)

	for rows.Next() {
		var shortID, originalURL string
		err = rows.Scan(&shortID, &originalURL)
		if err != nil {
			log.Warn().Err(err).Msg("GetData Cannot scan row")
			continue
		}
		data[shortID] = originalURL
	}

	return data, nil
}

// generateDollarSigns - генерирует строку вида "$1, $2, $3, ...",
// для использования в выражении IN запроса к базе данных.
// Параметры:
// - n - количество знаков доллара
// - start - начальное значение после знака доллара
// Возвращает строку с запятыми и знаками доллара.
func generateDollarSigns(n int, start int) string {
	result := "("
	for i := 0; i < n; i++ {
		if i == 0 {
			result += fmt.Sprintf("$%d", start)
		} else {
			result += fmt.Sprintf(", $%d", start+i)
		}
	}
	return result + ")"
}

// DeleteKeys - помечает записи в базе данных как удаленные
// добавляя префикс "-" к short_id.
// Параметры:
// - ctx - контекст
// - userID - идентификатор пользователя
// - keys - массив ключей
// Возвращает ошибку, если удаление не удалось.
func DeleteKeys(ctx context.Context, userID string, keys []any) error {
	if !IsConnected() {
		return errors.New("DeleteKeys. No connection to DB")
	}
	// если ключи не переданы, возвращаем успех
	if len(keys) == 0 {
		return nil
	}

	// готовим запрос для обновления записей
	query := `UPDATE urls 
	SET short_id = '-' || short_id 
	WHERE original_url LIKE $1 || '@%' 
	AND short_id IN ` + generateDollarSigns(len(keys), 2)

	// готовим аргументы для запроса
	args := make([]any, 0, len(keys))
	args = append(args, userID)
	args = append(args, keys...)

	// выполняем запрос
	// // Check if ctx canceled
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	// _, err := DB.ExecContext(ctx, query, args...)
	_, err := DB.Exec(query, args...)
	// log.Info().Msgf("*********************************")
	return err
}
