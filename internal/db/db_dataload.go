// Description: Функции для загрузки данных из базы данных в storage.

package db

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// LoadDataToStorage - загружает данные из базы данных в storage.
// Параметры:
// - ctx - контекст
// Возвращает ошибку, если загрузка данных не удалась.
func LoadDataToStorage(ctx context.Context) (err error) {
	if !IsConnected() {
		err = errors.New("LoadDataToStorage(). No connection to DB")
		log.Error().Err(err).Msg("LoadDataToStorage(). Cannot load data from DB")
		return err
	}
	data, err := GetData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("loadDataFromDB(). Cannot get data from DB")
		return err
	}
	storage.LoadData(data)
	log.Info().Msgf("%d Records loaded from database", len(data))
	return nil
}
