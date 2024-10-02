package main

import (
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
	db.Connect(1)
	db.MigrateUp("./migrations")
	server.ServeChi(config.Params.ServerAddress)
}
