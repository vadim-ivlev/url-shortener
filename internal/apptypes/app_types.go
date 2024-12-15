package apptypes

// Структура для хранения данных в памяти.
type URLShortener struct {
	Idx         int64  `json:"idx" db:"idx"`
	ShortID     string `json:"short_id" db:"short_id"`
	OriginalURL string `json:"original_url" db:"original_url"`
	UserID      string `json:"user_id" db:"user_id"`
	// Запись удалена. 0 - не удалена, 1 - удалена.
	Deleted int64 `json:"deleted" db:"deleted"`
}

// MemStoreInterface - интерфейс для работы с хранилищем записей URLShortener.
type MemStoreInterface interface {
	// Clear очищает хранилище.
	Clear() (err error)
	// AddRecord добавляет запись в хранилище.
	AddRecord(record URLShortener) (addedRecord URLShortener, created bool, err error)
	// AddRecords добавляет несколько записей в хранилище.
	AddRecords(records []URLShortener) (numAdded int, errs []error)
	// GetRecordByShortID извлекает запись по её shortID.
	GetRecordByShortID(shortID string) (record *URLShortener, err error)
	// GetRecordsByUserID извлекает все записи для пользователя.
	GetRecordsByUserID(userID string) (records []URLShortener, err error)
	// DeleteRecords удаляет несколько записей пользователя по их shortID.
	DeleteRecords(userID string, shortIDs []any) error
	// GetRecords возвращает все записи.
	GetRecords() (records []URLShortener, err error)
	// PrintRecords выводит содержимое нескольких записей.
	PrintRecords(limit int)
}

// ShortIDKeyFunc - функция для вычисления ключа для индекса по shortID.
//
// Параметры:
// - record - запись для вычисления ключа.
//
// Возвращает:
// - ключ для индекса.
func ShortIDKeyFunc(record URLShortener) string {
	return record.ShortID
}

// UserIDOriginalURLKeyFunc - функция для вычисления ключа для индекса по UserID+originalURL.
//
// Параметры:
// - record - запись для вычисления ключа.
//
// Возвращает:
// - ключ для индекса.
func UserIDOriginalURLKeyFunc(record URLShortener) string {
	return record.UserID + "@" + record.OriginalURL
}
