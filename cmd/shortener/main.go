package main

import (
	"github.com/vadim-ivlev/url-shortener/internal/app"
	"github.com/vadim-ivlev/url-shortener/internal/server"
)

func main() {
	// Инициализировать приложение
	app.InitApp()

	// Загрузить данные в хранилище
	app.LoadDataToMemStore()

	// Запустить сервер
	server.ServeChi()
}
