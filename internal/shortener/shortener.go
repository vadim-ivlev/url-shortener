package shortener

import (
	"fmt"
	"hash/fnv"
)

// Shorten генерирует укороченный ключ для данного значения.
func Shorten(value string) (key string) {
	hash := fnv.New32()
	hash.Write([]byte(value))
	return fmt.Sprintf("%X", hash.Sum32())
}
