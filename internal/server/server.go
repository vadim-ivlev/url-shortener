package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadim-ivlev/url-shortener/internal/handlers"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
)

// Using Chi router
func ServeChi(address string) {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Post("/", handlers.ShortenURLHandler)
	r.Get("/{id}", handlers.RedirectHandler)
	r.Route("/api", func(r chi.Router) {
		r.Use(contentTypeJSON)
		r.Post("/shorten", handlers.APIShortenHandler)
	})

	err := http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}
