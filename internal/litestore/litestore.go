package litestore

import (
	"errors"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
)

// Проверка на соответствие интерфейсу
var _ apptypes.MemStoreInterface = (*dbstore)(nil)

// DSN - строка подключения к базе данных
// var DSN = ":memory:"
var DSN = "url-shortener.db"

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
	d.dbPool, err = sqlx.Connect("sqlite", DSN)
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
func (d *dbstore) AddRecord(record apptypes.URLShortener) (addedRecord apptypes.URLShortener, created bool, err error) {
	// Проверяем нужно ли сохранять запись в файловое хранилище
	if !config.UseDatabase() {
		return addedRecord, false, nil
	}

	if err = d.IsConnected(); err != nil {
		return addedRecord, false, err
	}

	// Проверяем, есть ли уже такая запись в хранилище по UserID+originalURL
	err = d.dbPool.Get(&addedRecord, "SELECT idx, short_id, original_url, user_id, deleted FROM urls WHERE user_id = $1 AND original_url = $2", record.UserID, record.OriginalURL)
	if err == nil {
		return addedRecord, false, nil
	}

	// Если record.ShortID пустой, то генерируем новый
	if record.ShortID == "" {
		record.ShortID = shortener.Shorten(apptypes.UserIDOriginalURLKeyFunc(record))
	}

	// Добавляем запись в хранилище
	_, err = d.dbPool.Exec("INSERT INTO urls (idx, short_id, original_url, user_id, deleted) VALUES ($1, $2, $3, $4, $5)", record.Idx, record.ShortID, record.OriginalURL, record.UserID, record.Deleted)
	if err == nil {
		created = true
		addedRecord = record
	}

	return addedRecord, created, err
}

// AddRecords добавляет массив записей в хранилище.
//
// Параметры:
// - records - массив записей для добавления.
//
// Возвращает:
// - количество добавленных записей.
// - массив ошибок для записей которые не удалось добавить.
func (d *dbstore) AddRecords(records []apptypes.URLShortener) (numAdded int, errs []error) {
	for _, record := range records {
		_, _, err := d.AddRecord(record)
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
// - shortID - короткий идентификатор
//
// Возвращает запись apptypes.URLShortener и ошибку.
func (d *dbstore) GetRecordByShortID(shortID string) (record *apptypes.URLShortener, err error) {
	if err := d.IsConnected(); err != nil {
		return record, err
	}

	record = &apptypes.URLShortener{}

	err = d.dbPool.Get(record, "SELECT idx, short_id, original_url, user_id, deleted FROM urls WHERE short_id = $1", shortID)
	return record, err
}

// GetRecordsByUserID возвращает все записи пользователя.
//
// Параметры:
// - userID - идентификатор пользователя.
//
// Возвращает:
// - массив записей пользователя.
// - ошибку, если записи не найдены.
func (d *dbstore) GetRecordsByUserID(userID string) (records []apptypes.URLShortener, err error) {
	if err := d.IsConnected(); err != nil {
		return records, err
	}

	err = d.dbPool.Select(&records, "SELECT idx, short_id, original_url, user_id, deleted FROM urls WHERE user_id = $1", userID)

	return records, err
}

// DeleteRecords - помечает записи в базе данных как удаленные
// добавляя префикс "-" к short_id.
// Параметры:
// - userID - идентификатор пользователя
// - keys - массив ключей
// Возвращает ошибку, если удаление не удалось.
func (d *dbstore) DeleteRecords(userID string, keys []any) (err error) {
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
// Возвращает массив apptypes.URLShortener и ошибку.
func (d *dbstore) GetRecords() (data []apptypes.URLShortener, err error) {
	if err := d.IsConnected(); err != nil {
		return nil, err
	}
	err = d.dbPool.Select(&data, "SELECT idx, short_id, original_url, user_id, deleted FROM urls")
	return data, err
}

// PrintRecords выводит содержимое хранилища в консоль.
// limit - количество элементов, которые будут выведены.
func (d *dbstore) PrintRecords(limit int) {
	if err := d.IsConnected(); err != nil {
		log.Error().Err(err).Msg("PrintRecords")
		return
	}

	numRecords := 0
	err := d.dbPool.Get(&numRecords, "SELECT COUNT(*) FROM urls")
	if err != nil {
		log.Error().Err(err).Msg("PrintRecords")
		return
	}

	fmt.Printf("Memstore contains %d records\n", numRecords)

	records := []apptypes.URLShortener{}
	err = d.dbPool.Select(&records, "SELECT idx, short_id, original_url, user_id, deleted FROM urls LIMIT $1", limit)
	if err != nil {
		log.Error().Err(err).Msg("PrintRecords")
		return
	}

	for i, record := range records {
		fmt.Printf("#%2d %v\n", i, record)
	}

}
