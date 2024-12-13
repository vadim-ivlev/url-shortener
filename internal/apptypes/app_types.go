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

// Store - интерфейс для работы с хранилищем записей URLShortener.
type Store interface {
	// AddRecord добавляет запись в хранилище.
	AddRecord(record URLShortener) (addedRecord URLShortener, created bool, err error)
	// AddRecords добавляет несколько записей в хранилище.
	AddRecords(records []URLShortener) (numAdded int, errs []error)
	// GetRecordByShortID извлекает запись по её shortID.
	GetRecordByShortID(shortID string) (record *URLShortener, err error)
	// GetRecordsByUserID извлекает все записи для пользователя.
	GetRecordsByUserID(userID string) (records []URLShortener)
	// DeleteRecords удаляет несколько записей пользователя по их shortID.
	DeleteRecords(userID string, shortIDs []any) error
	// PrintRecords выводит содержимое нескольких записей.
	PrintRecords(limit int)
}
