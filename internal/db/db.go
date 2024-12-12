package db

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

var initSQL = `
-- urls - хранит список уникальных URL и их коротких ключей
CREATE TABLE IF NOT EXISTS urls (
    idx INTEGER,                      -- Индекс записи в memstore
    short_id TEXT PRIMARY KEY,         -- Короткий ключ
    original_url TEXT NOT NULL,        -- Оригинальный URL
    user_id TEXT,                      -- Идентификатор пользователя
    deleted INTEGER DEFAULT 0,         -- Флаг удаления
    UNIQUE (user_id, original_url)
);
`

// db - пул соединений с базой данных
var db *sqlx.DB = nil

// Connect - устанавливает соединение с базой данных
func Connect() (err error) {
	// Проверяем нужно ли подключаться к базе данных
	if !config.UseDatabase() {
		return nil
	}
	Disconnect()
	db, err = sqlx.Connect("postgres", config.Params.DatabaseDSN)
	if err != nil {
		return err
	}
	// Выполняем инициализацию базы данных
	_, err = db.Exec(initSQL)
	return err
}

// Disconnect - закрывает соединение с базой данных
func Disconnect() {
	if db != nil {
		db.Close()
	}
	db = nil
}

// IsConnected - проверяет, установлено ли соединение с базой данных
// TODO: should return error
func IsConnected() bool {
	return db != nil && db.Ping() == nil
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
	_, err := db.Exec("DELETE FROM urls")
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

	_, err := db.Exec("INSERT INTO urls (idx, short_id, original_url, user_id, deleted) VALUES ($1, $2, $3, $4, $5)", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted)
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

	_, err := db.Exec("UPDATE urls SET idx = $1,  short_id = $2, original_url = $3, user_id = $4, deleted = $5 WHERE short_id = $6", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted, record.ShortID)
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

	err = db.GetContext(ctx, &record, "SELECT idx, short_id, original_url, user_id, deleted FROM urls WHERE short_id = $1", shortID)
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
	err = db.GetContext(ctx, &data, "SELECT idx, short_id, original_url, user_id, deleted FROM urls")
	return data, err
}

// DeleteShortIDs - помечает записи в базе данных как удаленные
// добавляя префикс "-" к short_id.
// Параметры:
// - ctx - контекст
// - userID - идентификатор пользователя
// - keys - массив ключей
// Возвращает ошибку, если удаление не удалось.
func DeleteShortIDs(ctx context.Context, userID string, keys []any) (err error) {
	if !config.UseDatabase() {
		return nil
	}

	if !IsConnected() {
		return errors.New("DeleteKeys. No connection to DB")
	}
	// если ключи не переданы, возвращаем успех
	if len(keys) == 0 {
		return nil
	}

	// готовим запрос для обновления записей. https://jmoiron.github.io/sqlx/
	query, args, err := sqlx.In("UPDATE urls SET deleted = 1 WHERE user_id = ? AND short_id IN (?)", userID, keys)
	if err != nil {
		return err
	}
	// rebinding query to adapt to the DB driver's bindvar type
	query = db.Rebind(query)
	// выполняем запрос
	_, err = db.Exec(query, args...)

	return err
}
