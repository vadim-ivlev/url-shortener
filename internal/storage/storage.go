package storage

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

// dm —  экземпляр DoubleMap.
var dm *DoubleMap

type UrlsRec struct {
	ShortID     string
	OriginalURL string
	UserID      string
}

// DoubleMap - двухсторонняя карта для хранения отображения между оригинальными значениями и их укороченными ключами.
// valueToKey — это карта для хранения отображения от оригинальных значений к их укороченным ключам.
// keyToValue — это карта для хранения отображения от укороченных ключей к их оригинальным значениям.
// mu — это мьютекс для обеспечения потокобезопасных операций с картами.
// Эта реализация должна обеспечивать временную сложность O(1) для  операций Set и Get.
type DoubleMap struct {
	valueToKey map[string]string
	keyToValue map[string]string
	mutex      sync.Mutex
}

// Create порождает новый экземпляр DoubleMap.
func Create() {
	// Не был ли DoubleMap уже создан
	if dm != nil {
		return
	}
	dm = &DoubleMap{
		valueToKey: make(map[string]string),
		keyToValue: make(map[string]string),
	}
	log.Info().Msg("storage initialized")
}

// Clear очищает хранилище.
func Clear() {
	dm = &DoubleMap{
		valueToKey: make(map[string]string),
		keyToValue: make(map[string]string),
	}
	log.Info().Msg("storage cleared")
}

// Set сохраняет ключ и значени в DoubleMap.
// Сначала проверяется, существует ли значение уже в карте valueToKey. Если да, то возвращается существующий ключ.
// Если значение не существует, оно сохраняется с новым ключом и возвращается новый ключ.
// Новые отображения добавляются в обе карты.
// Возвращает ключ и флаг, указывающий, было ли новое значение добавлено в карту.
func Set(key, value string) (savedKey string, newKeyAdded bool) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Проверяем, существует ли уже укороченное значение
	if existingKey, exists := dm.valueToKey[value]; exists {
		return existingKey, false
	}

	// Сохраняем новое значение и ключ в обе карты
	dm.valueToKey[value] = key
	dm.keyToValue[key] = value

	return key, true
}

// LoadData - загружает данные из map[string]string, где ключ - short_id, значение - original_url, в storage.
func LoadData(data map[string]string) {
	for shortID, originalURL := range data {
		Set(shortID, originalURL)
	}
}

// Get возвращает значение для данного ключа.
// Если ключ не найден, возвращается пустая строка.
func Get(key string) (value string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Извлекаем значение
	value = dm.keyToValue[key]
	return value
}

// PrintContent выводит содержимое хранилища в консоль.
// limit - количество элементов, которые будут выведены.
func PrintContent(limit int) {
	log.Info().Msgf("RAM Storage contains %d records", len(dm.keyToValue))
	n := 0
	for k, v := range dm.keyToValue {
		n++
		if n > limit {
			break
		}
		fmt.Printf("%4v %v %v\n", n, k, v)
	}
}

// GetData - возвращает данные  в виде map[string]string,
// где ключ - short_id, значение - original_url.
func GetData() (data map[string]string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	data = make(map[string]string)
	for k, v := range dm.keyToValue {
		data[k] = v
	}
	return
}
