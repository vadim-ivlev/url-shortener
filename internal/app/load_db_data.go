// Description: Функции для загрузки данных из базы данных в storage.

package app

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/memstore"
)

// LoadDBDataToStorage - загружает данные из базы данных в storage.
// Параметры:
// - ctx - контекст
// Возвращает ошибку, если загрузка данных не удалась.
func LoadDBDataToStorage(ctx context.Context) (err error) {
	// Проверяем нужно ли загружать данные из базы данных
	if !config.UseDatabase() {
		return nil
	}

	// Проверяем, что есть соединение с базой данных
	if err := db.IsConnected(); err != nil {
		log.Error().Err(err).Msg("LoadDBDataToStorage(). Cannot load data from DB")
		return err
	}
	records, err := db.GetRecords(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("loadDataFromDB(). Cannot get data from DB")
		return err
	}
	numAdded, errs := memstore.Store.AddRecords(records)
	if len(errs) > 0 {
		log.Error().Errs("errors", errs).Msg("loadDataFromDB(). Errors while adding records to storage")
	}

	log.Info().Msgf("%d Records loaded from database", numAdded)
	return nil
}
