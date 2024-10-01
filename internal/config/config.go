package config

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v11"
)

// config - структура для хранения параметров приложения
type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

// Params - переменная для хранения параметров приложения
var Params config = config{}

func ParseCommandLine() {
	// Читаем параметры командной строки с значениями по умолчанию
	flag.StringVar(&Params.ServerAddress, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&Params.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&Params.FileStoragePath, "f", "./data/file-storage.txt", "File storage path")
	flag.Parse()

	// Читаем переменные окружения
	envVars := config{}
	if err := env.Parse(&envVars); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Если переменные окружения заданы, то используем их
	if envVars.ServerAddress != "" {
		Params.ServerAddress = envVars.ServerAddress
	}
	if envVars.BaseURL != "" {
		Params.BaseURL = envVars.BaseURL
	}
	if envVars.FileStoragePath != "" {
		Params.FileStoragePath = envVars.FileStoragePath
	}

	log.Info().Msg("Server Address: " + Params.ServerAddress)
	log.Info().Msg("Shortened Base URL: " + Params.BaseURL)
	log.Info().Msg("File Storage Path: " + Params.FileStoragePath)
}
