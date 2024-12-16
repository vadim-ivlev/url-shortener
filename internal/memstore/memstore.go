package memstore

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/arraystore"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/litestore"
	"github.com/vadim-ivlev/url-shortener/internal/pg"
)

// Проверка на соответствие интерфейсу
var _ apptypes.MemStoreInterface = (*memstore)(nil)

type memstore struct {
	apptypes.MemStoreInterface
}

var actualStore apptypes.MemStoreInterface

// New создает новое хранилище Urls.
func New() (st *memstore) {

	if config.Params.MemStore == "array" {
		actualStore = arraystore.New()
	} else {
		actualStore = litestore.New()
	}

	return &memstore{}
}

// Clear очищает хранилище.
func (m *memstore) Clear() error {
	// очищаем файловое хранилище
	filestorage.Clear()
	// очищаем базу данных
	pg.PGStore.Clear()
	// очищаем хранилище
	return actualStore.Clear()
}

// AddRecord добавляет запись в хранилище.
func (m *memstore) AddRecord(record apptypes.URLShortener) (addedRecord apptypes.URLShortener, created bool, err error) {
	addedRecord, created, err = actualStore.AddRecord(record)
	if err == nil && created {
		// Сохраняем запись в файловое хранилище
		err0 := filestorage.AddRecord(addedRecord)
		if err0 != nil {
			log.Error().Err(err0).Msg("Add() filestorage.AddRecord")
		}

		// Сохраняем в базу данных
		err1 := pg.PGStore.AddRecord(addedRecord)
		if err1 != nil {
			log.Error().Err(err1).Msg("Add() AddRecord")
		}
	}
	return
}

// AddRecords добавляет несколько записей в хранилище.
func (m *memstore) AddRecords(records []apptypes.URLShortener) (numAdded int, errs []error) {
	numAdded, errs = actualStore.AddRecords(records)
	log.Info().Msgf("AddRecords() numAdded: %d", numAdded)

	// Сохраняем записи в файловое хранилище
	recs, err0 := actualStore.GetRecords()
	if err0 != nil {
		log.Error().Err(err0).Msg("AddRecords() actualStore.GetRecords")
	} else {
		err00 := filestorage.DumpRecords(recs)
		if err00 != nil {
			log.Error().Err(err00).Msg("AddRecords() filestorage.SaveRecords")
		}
	}

	// Сохраняем записи в базу данных
	numAdded1, errs1 := pg.PGStore.AddRecords(records)
	if len(errs1) > 0 {
		log.Error().Msgf("AddRecords() db.PGStore.AddRecords: %v", errs1)
	}
	log.Info().Msgf("AddRecords() numAdded1: %d", numAdded1)

	return
}

// GetRecordByShortID извлекает запись по её shortID.
func (m *memstore) GetRecordByShortID(shortID string) (*apptypes.URLShortener, error) {
	return actualStore.GetRecordByShortID(shortID)
}

// GetRecordsByUserID извлекает все записи для пользователя.
func (m *memstore) GetRecordsByUserID(userID string) ([]apptypes.URLShortener, error) {
	return actualStore.GetRecordsByUserID(userID)
}

// DeleteRecords удаляет несколько записей пользователя по их shortID.
func (m *memstore) DeleteRecords(userID string, shortIDs []interface{}) (err error) {
	err = actualStore.DeleteRecords(userID, shortIDs)

	// получаем записи из хранилища
	records, err0 := actualStore.GetRecords()
	if err0 != nil {
		log.Error().Err(err0).Msg("DeleteRecords() actualStore.GetRecords")
	} else {
		// сохраняем записи в файловое хранилище
		err00 := filestorage.DumpRecords(records)
		if err00 != nil {
			log.Error().Err(err00).Msg("DeleteRecords() filestorage.SaveRecords")
		}
	}

	// удаляем записи из базы данных
	err1 := pg.PGStore.DeleteRecords(context.Background(), userID, shortIDs)
	if err1 != nil {
		log.Error().Err(err1).Msg("DeleteRecords() db.PGStore.DeleteRecords")
	}

	return err
}

// GetRecords возвращает все записи.
func (m *memstore) GetRecords() ([]apptypes.URLShortener, error) {
	return actualStore.GetRecords()
}

// PrintRecords выводит содержимое нескольких записей.
func (m *memstore) PrintRecords(limit int) {
	actualStore.PrintRecords(limit)
}
