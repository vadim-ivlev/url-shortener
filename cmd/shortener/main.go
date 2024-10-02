package main

import (
	// "log"
	"github.com/rs/zerolog/log"

	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/server"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func main() {
	logger.InitializeLogger()
	config.ParseCommandLine()
	storage.Create()
	_, err := filestorage.LoadData(config.Params.FileStoragePath)
	if err != nil {
		log.Warn().Err(err).Msg("Filestorage not found. Probably this is the first launch.")
	}
	storage.PrintKeyValue()
	server.ServeChi(config.Params.ServerAddress)
}
