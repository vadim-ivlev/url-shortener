// Description: Функции для загрузки данных из базы данных в storage.

package db

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// GetData - возвращает данные из базы данных в виде map[string]string,
// где ключ - short_id, значение - original_url.
func GetData() (data map[string]string, err error) {
	if !IsConnected() {
		return nil, errors.New("GetData. No connection to DB")
	}

	rows, err := DB.Queryx("SELECT short_id, original_url FROM urls")
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

	return data, err
}

// LoadData - загружает данные из базы данных в storage.
func LoadData() {
	if !IsConnected() {
		log.Error().Msg("loadDataFromDB(). No connection to DB")
		return
	}
	data, err := GetData()
	if err != nil {
		log.Warn().Err(err).Msg("loadDataFromDB(). Cannot get data from DB")
	}
	storage.LoadData(data)
	log.Info().Msg(fmt.Sprintf("%d Records loaded from database", len(data)))
}
