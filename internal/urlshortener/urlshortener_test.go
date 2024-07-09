package urlshortener

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func assertEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("expected (%+v) is not equal to actual (%+v)", expected, actual)
}

func Test_Urlshortener(t *testing.T) {

	key1 := Us.Shorten("example.com")
	fmt.Println("Укороченный ключ для 'example.com':", key1)

	expandedValue := Us.Expand(key1)
	fmt.Println("Расширенное значение для", key1, ":", expandedValue)

	assertEqual(t, "example.com", expandedValue)

	key2 := Us.Shorten("another.com")
	fmt.Println("Укороченный ключ для 'another.com':", key2)

	expandedValue2 := Us.Expand(key2)
	fmt.Println("Расширенное значение для", key2, ":", expandedValue2)

	assertEqual(t, "another.com", expandedValue2)

	expandedValue3 := Us.Expand("non-existent-key")
	fmt.Println("Расширенное значение для 'non-existent-key':", expandedValue3)

	assertEqual(t, "", expandedValue3)

}
