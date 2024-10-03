package main

import (
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/server"
)

func main() {
	// Инициализировать приложение
	app.InitApp()

	// Запустить сервер
	server.ServeChi()
}
