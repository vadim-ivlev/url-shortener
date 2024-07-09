package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadim-ivlev/url-shortener/internal/handlers"
)

// Using Chi router
func ServeChi(address string) {
	r := chi.NewRouter()

	r.Post("/", handlers.ShortenURLHandler)
	r.Get("/{id}", handlers.RedirectHandler)

	err := http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}
