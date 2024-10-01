package storage

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

// dm —  экземпляр DoubleMap.
var dm *DoubleMap

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
	log.Info().Msg("storage package initialized")
}

// Set сохраняет ключ и значени в DoubleMap.
// Сначала проверяется, существует ли значение уже в карте valueToKey. Если да, то возвращается существующий ключ.
// Если значение не существует, оно сохраняется с новым ключом и возвращается новый ключ.
// Новые отображения добавляются в обе карты.
// Возвращает ключ и флаг, указывающий, было ли значение добавлено в карту.
func Set(key, value string) (string, bool) {
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

// Get возвращает значение для данного ключа.
// Если ключ не найден, возвращается пустая строка.
func Get(key string) (value string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Извлекаем значение
	value = dm.keyToValue[key]
	return value
}

func PrintKeyValue() {
	fmt.Println("Storage: # key value >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	n := 1
	for k, v := range dm.keyToValue {
		fmt.Printf("%4v %v %v\n", n, k, v)
		n++
	}
	fmt.Println("Storage: <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
