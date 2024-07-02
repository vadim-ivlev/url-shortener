package server

import "net/http"

// addRouts добавляет обработчики маршрутов к мультиплексору.
func addRouts(mux *http.ServeMux) {
	// Зарегистрировать обработчики
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ShortenURLHandler(w, r)
		} else if r.Method == http.MethodGet {
			RedirectHandler(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
		}
	})
}

func Serve() {
	// Создать мультиплексор
	mux := http.NewServeMux()

	// Зарегистрировать маршруты и обработчики
	addRouts(mux)

	// Запустить сервер на порту 8080
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
