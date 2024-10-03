package main

import (
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/db"
	"github.com/vadim-ivlev/url-shortener/internal/logger"
	"github.com/vadim-ivlev/url-shortener/internal/server"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func main() {
	logger.InitializeLogger()
	config.ParseCommandLine()
	config.PrintParams()

	storage.Create()
	db.Connect(1)
	db.MigrateUp("./migrations")

	app.LoadDataToStorage()
	storage.PrintContent(0)

	server.ServeChi()
}
