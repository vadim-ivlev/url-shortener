package main

import (
	"github.com/vadim-ivlev/url-shortener/internal/config"
	"github.com/vadim-ivlev/url-shortener/internal/server"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func main() {
	config.ParseCommandLine()
	storage.Create()
	server.ServeChi(config.Address)
}
