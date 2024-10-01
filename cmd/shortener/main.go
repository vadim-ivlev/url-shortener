package main

import (
	// "log"
	"github.com/rs/zerolog/log"

	"github.com/vadim-ivlev/url-shortener/internal/config"
	filestorage "github.com/vadim-ivlev/url-shortener/internal/iter9-filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/server"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func main() {
	logger.InitializeLogger()
	config.ParseCommandLine()
	storage.Create()
	_, err := filestorage.LoadData(config.FileStoragePath)
	if err != nil {
		log.Warn().Err(err).Msg("Filestorage not found. Probably this is the first launch.")
	}
	server.ServeChi(config.Address)
}
