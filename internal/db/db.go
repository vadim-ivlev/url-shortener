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

// PGStore - хранилище записей URLShortener в базе данных.
var PGStore *dbstore

// dbstore - структура для хранения данных в базе данных.
type dbstore struct {
	// dbPool - пул соединений с базой данных
	dbPool *sqlx.DB
}

// New создает новое хранилище Urls.
func New() *dbstore {
	return &dbstore{}
}

// Connect - устанавливает соединение с базой данных
func (d *dbstore) Connect() (err error) {
	// Проверяем нужно ли подключаться к базе данных
	if !config.UseDatabase() {
		return nil
	}
	d.Disconnect()
	d.dbPool, err = sqlx.Connect("postgres", config.Params.DatabaseDSN)
	if err != nil {
		return err
	}
	// Выполняем инициализацию базы данных
	_, err = d.dbPool.Exec(initSQL)
	return err
}

// Disconnect - закрывает соединение с базой данных
func (d *dbstore) Disconnect() {
	if d.dbPool != nil {
		d.dbPool.Close()
	}
	d.dbPool = nil
}

// IsConnected - проверяет, установлено ли соединение с базой данных
func (d *dbstore) IsConnected() error {
	// return db != nil && db.Ping() == nil
	if d.dbPool == nil {
		return errors.New("no connection to DB")
	}
	return d.dbPool.Ping()
}

// Clear - очищает таблицу urls
//
// Возвращает ошибку, если очистка не удалась.
func (d *dbstore) Clear() error {
	if !config.UseDatabase() {
		return nil
	}

	if err := d.IsConnected(); err != nil {
		return err
	}
	_, err := d.dbPool.Exec("DELETE FROM urls")
	return err
}

// AddRecord - добавляет запись в базу данных.
//
// Параметры:
// - record - запись для сохранения.
//
// Возвращает ошибку, если запись не удалась.
func (d *dbstore) AddRecord(record apptypes.URLShortener) (err error) {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseDatabase() {
		return nil
	}

	if err = d.IsConnected(); err != nil {
		return err
	}

	_, err = d.dbPool.Exec("INSERT INTO urls (idx, short_id, original_url, user_id, deleted) VALUES ($1, $2, $3, $4, $5)", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted)
	return err
}

// AddRecords добавляет несколько записей в хранилище.
func (d *dbstore) AddRecords(records []apptypes.URLShortener) (numAdded int, errs []error) {
	for _, record := range records {
		err := d.AddRecord(record)
		if err != nil {
			errs = append(errs, err)
		} else {
			numAdded++
		}
	}
	return numAdded, errs
}

// UpdateRecord - обновляет запись в базе данных.
//
// Параметры:
// - record - запись для сохранения.
//
// Возвращает ошибку, если запись не удалась.
func (d *dbstore) UpdateRecord(record apptypes.URLShortener) error {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseDatabase() {
		return nil
	}

	if err := d.IsConnected(); err != nil {
		return err
	}

	_, err := d.dbPool.Exec("UPDATE urls SET idx = $1,  short_id = $2, original_url = $3, user_id = $4, deleted = $5 WHERE short_id = $6", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted, record.ShortID)
	return err
}

// GetRecordByShortID - возвращает запись из базы данных по short_id.
//
// Параметры:
// - ctx - контекст
// - shortID - короткий идентификатор
//
// Возвращает запись apptypes.URLShortener и ошибку.
func (d *dbstore) GetRecordByShortID(ctx context.Context, shortID string) (record apptypes.URLShortener, err error) {
	if err := d.IsConnected(); err != nil {
		return record, err
	}

	err = d.dbPool.GetContext(ctx, &record, "SELECT idx, short_id, original_url, user_id, deleted FROM urls WHERE short_id = $1", shortID)
	return record, err
}

// // GetRecordsByUserID извлекает все записи для пользователя.
// GetRecordsByUserID(userID string) (records []URLShortener)

// DeleteRecords - помечает записи в базе данных как удаленные
// добавляя префикс "-" к short_id.
// Параметры:
// - ctx - контекст
// - userID - идентификатор пользователя
// - keys - массив ключей
// Возвращает ошибку, если удаление не удалось.
func (d *dbstore) DeleteRecords(ctx context.Context, userID string, keys []any) (err error) {
	if !config.UseDatabase() {
		return nil
	}

	if err := d.IsConnected(); err != nil {
		return err
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
	query = d.dbPool.Rebind(query)
	// выполняем запрос
	_, err = d.dbPool.Exec(query, args...)

	return err
}

// GetRecords - возвращает данные из базы данных в виде массива apptypes.URLShortener.
//
// Параметры:
// - ctx - контекст
//
// Возвращает массив apptypes.URLShortener и ошибку.
func (d *dbstore) GetRecords(ctx context.Context) (data []apptypes.URLShortener, err error) {
	if err := d.IsConnected(); err != nil {
		return nil, err
	}
	err = d.dbPool.SelectContext(ctx, &data, "SELECT idx, short_id, original_url, user_id, deleted FROM urls")
	return data, err
}
