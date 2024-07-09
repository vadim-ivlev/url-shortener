package shortener

import (
	"fmt"
	"sync"
)

// counter счетчик для генерации уникальных ключей.
var counter int

// потокобезопасность
var mutex = &sync.Mutex{}

// Shorten генерирует укороченный ключ для данного значения.
func Shorten(value string) (key string) {
	mutex.Lock()
	counter++
	mutex.Unlock()
	return fmt.Sprintf("short%d", counter)
}
