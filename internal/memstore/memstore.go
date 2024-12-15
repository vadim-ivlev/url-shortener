package memstore

import (
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
)

// Проверка на соответствие интерфейсу
var _ apptypes.MemStoreInterface = (*memstore)(nil)

type memstore struct {
	apptypes.MemStoreInterface
}

var actualStore *memstore

// New создает новое хранилище Urls.
func New() (st *memstore) {

	// if config.Params.MemStore == "array" {
	// 	actualStore = arraystore.New()
	// } else {
	// 	actualStore = litestore.New()
	// }

	actualStore = &memstore{}

	return actualStore
}

// Clear очищает хранилище.
func (m *memstore) Clear() error {
	return nil
}

// AddRecord добавляет запись в хранилище.
func (m *memstore) AddRecord(record apptypes.URLShortener) (apptypes.URLShortener, bool, error) {
	return record, true, nil
}

// AddRecords добавляет несколько записей в хранилище.
func (m *memstore) AddRecords(records []apptypes.URLShortener) (int, []error) {
	return len(records), nil
}

// GetRecordByShortID извлекает запись по её shortID.
func (m *memstore) GetRecordByShortID(shortID string) (*apptypes.URLShortener, error) {
	return nil, nil
}

// GetRecordsByUserID извлекает все записи для пользователя.
func (m *memstore) GetRecordsByUserID(userID string) ([]apptypes.URLShortener, error) {
	return nil, nil
}

// DeleteRecords удаляет несколько записей пользователя по их shortID.
func (m *memstore) DeleteRecords(userID string, shortIDs []interface{}) error {
	return nil
}

// GetRecords возвращает все записи.
func (m *memstore) GetRecords() ([]apptypes.URLShortener, error) {
	return nil, nil
}

// PrintRecords выводит содержимое нескольких записей.
func (m *memstore) PrintRecords(limit int) {
}
