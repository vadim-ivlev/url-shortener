package main

import (
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
)

const (
	addr    = ":8080"  // адрес сервера
	maxSize = 10000000 // будем растить слайс до 10 миллионов элементов
)

func foo() {
	// полезная нагрузка
	for {
		s := make([]int, 0, maxSize)
		for i := 0; i < maxSize; i++ {
			s = append(s, i)
			_ = s
		}
	}
}

func main() {
	go foo()                       // запускаем полезную нагрузку в фоне
	http.ListenAndServe(addr, nil) // запускаем сервер
}
