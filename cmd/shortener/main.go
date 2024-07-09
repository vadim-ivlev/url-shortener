package main

import (
	"fmt"

	"github.com/vadim-ivlev/url-shortener/internal/server"
	"github.com/vadim-ivlev/url-shortener/internal/storage"
)

func main() {
	storage.Create()
	fmt.Println("Server starting at http://localhost:8080/")
	server.Serve()
}
