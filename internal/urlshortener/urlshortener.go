package urlshortener

import (
	"fmt"
	"sync"
)

// Us — это глобальная переменная, которая представляет экземпляр Urlshortener.
var Us *Urlshortener

// Urlshortener
// valueToKey — это карта для хранения отображения от оригинальных значений к их укороченным ключам.
// keyToValue — это карта для хранения отображения от укороченных ключей к их оригинальным значениям.
// mu — это мьютекс для обеспечения потокобезопасных операций с картами.
// counter используется для генерации уникальных ключей.
// Эта реализация должна обеспечивать временную сложность O(1) для обеих операций Shorten и Expand.
type Urlshortener struct {
	valueToKey map[string]string
	keyToValue map[string]string
	mu         sync.Mutex
	counter    int
}

// NewUrlshortener инициализирует новый экземпляр Urlshortener.
func NewUrlshortener() *Urlshortener {
	return &Urlshortener{
		valueToKey: make(map[string]string),
		keyToValue: make(map[string]string),
		counter:    0,
	}
}

// Shorten генерирует укороченный ключ для данного значения.
// Сначала проверяется, существует ли значение уже в карте valueToKey. Если да, то возвращается существующий ключ.
// Если значение не существует, генерируется новый ключ с использованием текущего значения счетчика.
// Новые отображения добавляются в обе карты, и значение счетчика увеличивается.
func (s *Urlshortener) Shorten(value string) (key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли уже укороченное значение
	if key, exists := s.valueToKey[value]; exists {
		return key
	}

	// Генерируем новый ключ
	key = fmt.Sprintf("short%d", s.counter)
	s.counter++

	// Сохраняем обе маппинги
	s.valueToKey[value] = key
	s.keyToValue[key] = value

	return key
}

// Expand возвращает оригинальное значение для данного ключа.
func (s *Urlshortener) Expand(key string) (value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Извлекаем оригинальное значение
	value = s.keyToValue[key]
	return value
}

func init() {
	fmt.Println("urlshortener package initialized")
	Us = NewUrlshortener()
}
