package arraystore

import (
	"fmt"

	"sync"

	"errors"

	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
)

// Проверка на соответствие интерфейсу
var _ apptypes.MemStoreInterface = (*mstore)(nil)

var ErrRecordNotFound = errors.New("record not found")

// mstore - структура для хранения данных в памяти.
// Содержит массив записей UrlShortener.
// Потокобезопасна.
// Имеет два уникальных индекса для быстрого поиска записей по shortID и UserID+originalURL.
type mstore struct {
	// Записи
	records []apptypes.URLShortener
	// Мьютекс для защиты записей
	mutex sync.Mutex
	// Индекс для поиска записей по shortID
	idxShortID *Index
	// Индекс для поиска записей по UserID+originalURL
	idxUserIDOriginalURL *Index
}

// New создает новое хранилище Urls.
func New() *mstore {
	return &mstore{
		records:              make([]apptypes.URLShortener, 0),
		idxShortID:           NewIndex(apptypes.ShortIDKeyFunc),
		idxUserIDOriginalURL: NewIndex(apptypes.UserIDOriginalURLKeyFunc),
	}
}

// Clear очищает хранилище.
func (u *mstore) Clear() (err error) {
	// очищаем файловое хранилище
	filestorage.Clear()
	// очищаем базу данных
	db.PGStore.Clear()

	// очищаем хранилище
	u.records = make([]apptypes.URLShortener, 0)
	// пересоздаем индексы
	u.idxShortID = NewIndex(apptypes.ShortIDKeyFunc)
	u.idxUserIDOriginalURL = NewIndex(apptypes.UserIDOriginalURLKeyFunc)
	return err
}

// AddRecord добавляет запись в хранилище.
//
// Параметры:
// - record - запись для добавления.
//
// Возвращает:
// - добавленную запись, или ту, что уже есть в хранилище.
// - true, если запись была добавлена, false, если запись уже есть в хранилище.
// - ошибку, если запись не удалось добавить.
// TODO: get rid of methods
func (u *mstore) AddRecord(record apptypes.URLShortener) (addedRecord apptypes.URLShortener, created bool, err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Проверяем, есть ли уже такая запись в хранилище по UserID+originalURL
	if idx, ok := u.idxUserIDOriginalURL.Get(record); ok {
		return u.records[idx], false, nil
	}

	// Если record.ShortID пустой, то генерируем новый
	if record.ShortID == "" {
		record.ShortID = shortener.Shorten(apptypes.UserIDOriginalURLKeyFunc(record))
	}
	// Вычисляем idx
	record.Idx = int64(len(u.records))

	// Добавляем запись в хранилище
	u.records = append(u.records, record)

	// Добавляем запись в индексы
	u.idxShortID.Add(record, record.Idx)
	u.idxUserIDOriginalURL.Add(record, record.Idx)

	// Сохраняем запись в файловое хранилище
	err0 := filestorage.AddRecord(record)
	if err0 != nil {
		log.Error().Err(err0).Msg("Add() filestorage.AddRecord")
	}

	// Сохраняем в базу данных
	err1 := db.PGStore.AddRecord(record)
	if err1 != nil {
		log.Error().Err(err1).Msg("Add() AddRecord")
	}

	return record, true, nil
}

// AddRecords добавляет массив записей в хранилище.
//
// Параметры:
// - records - массив записей для добавления.
//
// Возвращает:
// - количество добавленных записей.
// - массив ошибок для записей которые не удалось добавить.
func (u *mstore) AddRecords(records []apptypes.URLShortener) (numAdded int, errs []error) {
	for _, record := range records {
		_, _, err := u.AddRecord(record)
		if err != nil {
			errs = append(errs, err)
		} else {
			numAdded++
		}
	}
	return numAdded, errs
}

// GetRecordByShortID возвращает запись по shortID.
//
// Параметры:
// - shortID - shortID записи.
//
// Возвращает:
// - запись, если она найдена или nil
// - ошибку, если запись не найдена.
func (u *mstore) GetRecordByShortID(shortID string) (record *apptypes.URLShortener, err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	idx, ok := u.idxShortID.Get(apptypes.URLShortener{ShortID: shortID})
	if !ok {
		return nil, ErrRecordNotFound
	}
	return &u.records[idx], nil
}

// GetRecordsByUserID возвращает все записи пользователя.
//
// Параметры:
// - userID - идентификатор пользователя.
//
// Возвращает:
// - массив записей пользователя.
// - ошибку.
func (u *mstore) GetRecordsByUserID(userID string) (records []apptypes.URLShortener, err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	result := make([]apptypes.URLShortener, 0)
	for _, record := range u.records {
		if record.UserID == userID {
			result = append(result, record)
		}
	}

	return result, nil
}

// delete - делает пометку записи как удаленную.
// Удалить ключ может только пользователь его создавший.
//
// Параметры:
// - userID - идентификатор пользователя
// - key - ключ
//
// Возвращает ошибку
func (u *mstore) delete(userID, shortID string) error {
	// Блокируем доступ к хранилищу
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Проверяем, существует ли запись
	idx, ok := u.idxShortID.Get(apptypes.URLShortener{ShortID: shortID})
	if !ok {
		return ErrRecordNotFound
	}

	// Проверяем, что пользователь удаляет свой ключ
	if u.records[idx].UserID != userID {
		return fmt.Errorf("ShortID %s  does not belong to user %s", shortID, userID)
	}

	// Проверяем, не был ли ключ уже удален
	if u.records[idx].Deleted != 0 {
		return fmt.Errorf("ShortID %s already deleted", shortID)
	}

	// Помечаем ключ как удаленный
	u.records[idx].Deleted = 1

	// Сохраняем запись в файловое хранилище
	// TODO: too many writes
	filestorage.DumpRecords(u.records)

	// Сохраняем в базу данных
	err1 := db.PGStore.UpdateRecord(u.records[idx])
	// log.Info().Msgf("DEL >>> Record %#v deleted", u.Records[idx])
	if err1 != nil {
		log.Error().Err(err1).Msg("DeleteKeysFromStore")
	}

	return nil
}

// DeleteRecords - удаляет множество записи из хранилища в горутинах.
//
// Параметры:
// - userID - идентификатор пользователя
// - keys - массив ShortID
//
// Возвращает ошибку
// TODO: OPTIMISE:
func (u *mstore) DeleteRecords(userID string, shortIDs []any) error {
	for _, shortID := range shortIDs {
		go func(shortID string) {
			err := u.delete(userID, shortID)
			if err != nil {
				fmt.Println(err)
			}
		}(shortID.(string))
	}

	return nil
}

// GetRecords возвращает все записи.
func (u *mstore) GetRecords() (records []apptypes.URLShortener, err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.records, nil
}

// PrintRecords выводит содержимое хранилища в консоль.
// limit - количество элементов, которые будут выведены.
func (u *mstore) PrintRecords(limit int) {
	log.Info().Msgf("Memstore contains %d records", len(u.records))
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if limit > len(u.records) {
		limit = len(u.records)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("%+v\n", u.records[i])
	}
}
