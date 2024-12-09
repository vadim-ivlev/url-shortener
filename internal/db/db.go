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

	_, err := DB.Exec("INSERT INTO urls (idx, short_id, original_url, user_id, deleted) VALUES ($1, $2, $3, $4, $5)", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted)
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

	_, err := DB.Exec("UPDATE urls SET idx = $1,  short_id = $2, original_url = $3, user_id = $4, deleted = $5 WHERE short_id = $6", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted, record.ShortID)
	return err
}

// Clear - очищает таблицу urls
//
// Возвращает ошибку, если очистка не удалась.
func Clear() error {
	if !config.UseDatabase() {
		return nil
	}

	if !IsConnected() {
		return errors.New("Clear. No connection to DB")
	}
	_, err := DB.Exec("DELETE FROM urls")
	return err
}

// GetByShortID - возвращает запись из базы данных по short_id.
//
// Параметры:
// - ctx - контекст
// - shortID - короткий идентификатор
//
// Возвращает запись apptypes.URLShortener и ошибку.
func GetByShortID(ctx context.Context, shortID string) (record apptypes.URLShortener, err error) {
	if !IsConnected() {
		return record, errors.New("GetByShortID. No connection to DB")
	}

	err = DB.GetContext(ctx, &record, "SELECT idx, short_id, original_url, user_id, deleted FROM urls WHERE short_id = $1", shortID)
	return record, err
}

// GetRecords - возвращает данные из базы данных в виде массива apptypes.URLShortener.
//
// Параметры:
// - ctx - контекст
//
// Возвращает массив apptypes.URLShortener и ошибку.
func GetRecords(ctx context.Context) (data []apptypes.URLShortener, err error) {
	if !IsConnected() {
		return nil, errors.New("GetData. No connection to DB")
	}

	rows, err := DB.QueryxContext(ctx, "SELECT idx, short_id, original_url, user_id, deleted FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	data = make([]apptypes.URLShortener, 0)

	for rows.Next() {
		var record apptypes.URLShortener
		err = rows.Scan(&record.Idx, &record.ShortID, &record.OriginalURL, &record.UserID, &record.Deleted)
		if err != nil {
			log.Warn().Err(err).Msg("GetData Cannot scan row")
			continue
		}
		data = append(data, record)
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
