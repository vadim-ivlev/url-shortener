package main

import (
	"fmt"
	"net/http"

	"github.com/vadim-ivlev/url-shortener/internal/server"
)

func main() {

	fmt.Println("Server starting at http://localhost:8080/")

	// Создать мультиплексор
	mux := http.NewServeMux()

	// Зарегистрировать обработчики
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			server.ShortenURLHandler(w, r)
		} else if r.Method == http.MethodGet {
			server.RedirectHandler(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
		}
	})

	// Запустить сервер на порту 8080
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
