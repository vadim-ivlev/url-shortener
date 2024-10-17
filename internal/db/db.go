package db

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
