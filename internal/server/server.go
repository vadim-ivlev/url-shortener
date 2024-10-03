package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/compression"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/handlers"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
)

// ServeChi запускает сервер на порту, указанном в конфигурации.
func ServeChi() {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	r.Use(compression.GzipMiddleware)
	r.Post("/", handlers.ShortenURLHandler)
	r.Get("/{id}", handlers.RedirectHandler)
	r.Get("/ping", handlers.PingHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(contentTypeJSON)
		r.Post("/shorten", handlers.APIShortenHandler)
		r.Post("/api/shorten/batch", handlers.APIShortenBatchHandler)
	})

	address := config.Params.ServerAddress
	log.Info().Str("address", address).Msg("Starting the server at the ...")
	err := http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}
