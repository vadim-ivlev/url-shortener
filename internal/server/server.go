package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadim-ivlev/url-shortener/internal/handlers"
)

// Using Chi router
func ServeChi() {
	r := chi.NewRouter()

	r.Post("/", handlers.ShortenURLHandler)
	r.Get("/{id}", handlers.RedirectHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
