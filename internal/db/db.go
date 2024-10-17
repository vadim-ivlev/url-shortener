package db

import (
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
// - shortID - укороченный ID.
// - originalURL - оригинальный URL.
// Возвращает ошибку, если запись не удалась.
func Store(shortID, originalURL string) error {
	if !IsConnected() {
		return errors.New("Store. No connection to DB")
	}
	_, err := DB.Exec("INSERT INTO urls (short_id, original_url) VALUES ($1, $2)", shortID, originalURL)
	return err
}

// Clear - очищает таблицу urls
func Clear() error {
	if !IsConnected() {
		return errors.New("Clear. No connection to DB")
	}
	_, err := DB.Exec("DELETE FROM urls")
	return err
}
