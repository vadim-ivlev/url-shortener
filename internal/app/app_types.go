package app

// Структура для хранения данных в памяти.
type UrlShortener struct {
	Idx         int64  `json:"idx" db:"idx"`
	ShortID     string `json:"short_id" db:"short_id"`
	OriginalURL string `json:"original_url" db:"original_url"`
	UserID      string `json:"user_id" db:"user_id"`
	// Запись удалена. 0 - не удалена, 1 - удалена.
	Deleted int64 `json:"deleted" db:"deleted"`
}
