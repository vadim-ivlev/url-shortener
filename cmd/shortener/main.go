package main

import (
	// "log"

	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/filestorage"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/server"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func main() {
	logger.InitializeLogger()
	config.ParseCommandLine()
	config.PrintParams()
	storage.Create()
	filestorage.LoadDataAndLog(config.Params.FileStoragePath)
	db.Connect(3)
	server.ServeChi(config.Params.ServerAddress)
}
