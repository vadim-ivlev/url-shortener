// Description: Функции для загрузки данных из базы данных в storage.

package app

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

// LoadDBDataToStorage - загружает данные из базы данных в storage.
// Параметры:
// - ctx - контекст
// Возвращает ошибку, если загрузка данных не удалась.
func LoadDBDataToStorage(ctx context.Context) (err error) {
	if !db.IsConnected() {
		err = errors.New("LoadDBDataToStorage(). No connection to DB")
		log.Error().Err(err).Msg("LoadDBDataToStorage(). Cannot load data from DB")
		return err
	}
	data, err := db.GetData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("loadDataFromDB(). Cannot get data from DB")
		return err
	}
	storage.LoadData(data)
	log.Info().Msgf("%d Records loaded from database", len(data))
	return nil
}
