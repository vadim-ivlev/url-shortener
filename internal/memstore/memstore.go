package memstore

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

// Store - хранилище Urls.
var Store *Urls = NewUrls()

var ErrRecordNotFound = errors.New("record not found")

// Urls - структура для хранения данных в памяти.
// Содержит массив записей UrlShortener.
// Потокобезопасна.
// Имеет два уникальных индекса для быстрого поиска записей по shortID и UserID+originalURL.
type Urls struct {
	// Записи
	Records []apptypes.URLShortener
	// Мьютекс для защиты записей
	mutex sync.Mutex
	// Индекс для поиска записей по shortID
	idxShortID *Index
	// Индекс для поиска записей по UserID+originalURL
	idxUserIDOriginalURL *Index
}

func idxShortIDKeyFunc(record apptypes.URLShortener) string {
	return record.ShortID
}

func idxUserIDOriginalURLKeyFunc(record apptypes.URLShortener) string {
	return record.UserID + "@" + record.OriginalURL
}

// NewUrls создает новое хранилище Urls.
func NewUrls() *Urls {
	return &Urls{
		Records:              make([]apptypes.URLShortener, 0),
		idxShortID:           NewIndex(idxShortIDKeyFunc),
		idxUserIDOriginalURL: NewIndex(idxUserIDOriginalURLKeyFunc),
	}
}

// Add добавляет запись в хранилище.
//
// Параметры:
// - record - запись для добавления.
//
// Возвращает:
// - добавленную запись, или ту, что уже есть в хранилище.
// - true, если запись была добавлена, false, если запись уже есть в хранилище.
// - ошибку, если запись не удалось добавить.
func (u *Urls) Add(record apptypes.URLShortener) (addedRecord apptypes.URLShortener, created bool, err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Проверяем, есть ли уже такая запись в хранилище
	if idx, ok := u.idxUserIDOriginalURL.Get(record); ok {
		return u.Records[idx], false, nil
	}

	// Если record.ShortID пустой, то генерируем новый
	if record.ShortID == "" {
		record.ShortID = shortener.Shorten(idxUserIDOriginalURLKeyFunc(record))
	}
	// Вычисляем idx
	record.Idx = int64(len(u.Records))

	// Добавляем запись в хранилище
	u.Records = append(u.Records, record)

	// Добавляем запись в индексы
	u.idxShortID.Add(record, record.Idx)
	u.idxUserIDOriginalURL.Add(record, record.Idx)

	// Сохраняем запись в файловое хранилище
	err0 := filestorage.AddRecord(record)
	if err0 != nil {
		log.Error().Err(err0).Msg("Add() filestorage.AddRecord")
	}

	// Сохраняем в базу данных
	err1 := db.AddRecord(record)
	if err1 != nil {
		log.Error().Err(err1).Msg("Add() AddRecord")
	}

	return record, true, nil
}

// GetByShortID возвращает запись по shortID.
//
// Параметры:
// - shortID - shortID записи.
//
// Возвращает:
// - запись, если она найдена или nil
// - ошибку, если запись не найдена.
func (u *Urls) GetByShortID(shortID string) (record *apptypes.URLShortener, err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	idx, ok := u.idxShortID.Get(apptypes.URLShortener{ShortID: shortID})
	if !ok {
		return nil, ErrRecordNotFound
	}
	return &u.Records[idx], nil
}

// GetByUserID возвращает все записи пользователя.
//
// Параметры:
// - userID - идентификатор пользователя.
//
// Возвращает:
// - массив записей пользователя.
func (u *Urls) GetByUserID(userID string) (records []apptypes.URLShortener) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	result := make([]apptypes.URLShortener, 0)
	for _, record := range u.Records {
		if record.UserID == userID {
			result = append(result, record)
		}
	}

	return result
}

// delete - делает пометку записи как удаленную.
// Удалить ключ может только пользователь его создавший.
//
// Параметры:
// - userID - идентификатор пользователя
// - key - ключ
//
// Возвращает ошибку
func (u *Urls) delete(userID, shortID string) error {
	// Блокируем доступ к хранилищу
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Проверяем, существует ли запись
	idx, ok := u.idxShortID.Get(apptypes.URLShortener{ShortID: shortID})
	if !ok {
		return ErrRecordNotFound
	}

	// Проверяем, что пользователь удаляет свой ключ
	if u.Records[idx].UserID != userID {
		return fmt.Errorf("ShortID %s  does not belong to user %s", shortID, userID)
	}

	// Проверяем, не был ли ключ уже удален
	if u.Records[idx].Deleted != 0 {
		return fmt.Errorf("ShortID %s already deleted", shortID)
	}

	// Помечаем ключ как удаленный
	u.Records[idx].Deleted = 1

	// // TODO:?  Удаляем запись из индексов
	// u.idxShortID.Delete(u.Records[idx])
	// u.idxUserIDOriginalURL.Delete(u.Records[idx])

	// Сохраняем запись в файловое хранилище
	filestorage.DumpRecords(u.Records)

	// Сохраняем в базу данных
	err1 := db.UpdateRecord(u.Records[idx])
	// log.Info().Msgf("DEL >>> Record %#v deleted", u.Records[idx])
	if err1 != nil {
		log.Error().Err(err1).Msg("DeleteKeysFromStore")
	}

	return nil
}

// DeleteShortIDs - удаляет множество записи из хранилища в горутинах.
//
// Параметры:
// - userID - идентификатор пользователя
// - keys - массив ShortID
//
// Возвращает ошибку
func (u *Urls) DeleteShortIDs(userID string, shortIDs []any) error {
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

// PrintContent выводит содержимое хранилища в консоль.
// limit - количество элементов, которые будут выведены.
func (u *Urls) PrintContent(limit int) {
	log.Info().Msgf("Memstore contains %d records", len(u.Records))
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if limit > len(u.Records) {
		limit = len(u.Records)
	}

	for i := 0; i < limit; i++ {
		// fmt.Println(config.JSONString(u.Records[i]))
		fmt.Printf("%+v\n", u.Records[i])
	}
}
