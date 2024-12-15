package memstore

import (
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/arraystore"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/litestore"
)

// Проверка на соответствие интерфейсу
var _ apptypes.MemStoreInterface = (*memstore)(nil)

type memstore struct {
	apptypes.MemStoreInterface
}

var actualStore apptypes.MemStoreInterface

// New создает новое хранилище Urls.
func New() (st *memstore) {

	if config.Params.MemStore == "array" {
		actualStore = arraystore.New()
	} else {
		actualStore = litestore.New()
	}

	return &memstore{}
}

// Clear очищает хранилище.
func (m *memstore) Clear() error {
	return actualStore.Clear()
}

// AddRecord добавляет запись в хранилище.
func (m *memstore) AddRecord(record apptypes.URLShortener) (apptypes.URLShortener, bool, error) {
	return actualStore.AddRecord(record)
}

// AddRecords добавляет несколько записей в хранилище.
func (m *memstore) AddRecords(records []apptypes.URLShortener) (int, []error) {
	return actualStore.AddRecords(records)
}

// GetRecordByShortID извлекает запись по её shortID.
func (m *memstore) GetRecordByShortID(shortID string) (*apptypes.URLShortener, error) {
	return actualStore.GetRecordByShortID(shortID)
}

// GetRecordsByUserID извлекает все записи для пользователя.
func (m *memstore) GetRecordsByUserID(userID string) ([]apptypes.URLShortener, error) {
	return actualStore.GetRecordsByUserID(userID)
}

// DeleteRecords удаляет несколько записей пользователя по их shortID.
func (m *memstore) DeleteRecords(userID string, shortIDs []interface{}) error {
	return actualStore.DeleteRecords(userID, shortIDs)
}

// GetRecords возвращает все записи.
func (m *memstore) GetRecords() ([]apptypes.URLShortener, error) {
	return actualStore.GetRecords()
}

// PrintRecords выводит содержимое нескольких записей.
func (m *memstore) PrintRecords(limit int) {
	actualStore.PrintRecords(limit)
}
