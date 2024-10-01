package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type envConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var Address string
var BaseURL string
var FileStoragePath string

func ParseCommandLine() {
	// Читаем параметры командной строки с значениями по умолчанию
	flag.StringVar(&Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&FileStoragePath, "f", "./data/file-storage.txt", "File storage path")
	flag.Parse()

	// Читаем переменные окружения
	cfg := envConfig{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Если переменные окружения заданы, то используем их
	if cfg.ServerAddress != "" {
		Address = cfg.ServerAddress
	}
	if cfg.BaseURL != "" {
		BaseURL = cfg.BaseURL
	}
	if cfg.FileStoragePath != "" {
		FileStoragePath = cfg.FileStoragePath
	}

	fmt.Println("Server Address:", Address)
	fmt.Println("Shortened Base URL:", BaseURL)
	fmt.Println("File Storage Path:", FileStoragePath)
}
