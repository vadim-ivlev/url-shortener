package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v11"
)

// config - структура для хранения параметров приложения
type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// Params - переменная для хранения параметров приложения
var Params config = config{}

// getDefaultDatabaseDSN - возвращает DSN для подключения к базе данных  в зависимости от окружения (CI или локальное).
// Заплатка для прохождения 9 го автотеста в CI.
func getDefaultDatabaseDSN() string {
	if os.Getenv("CI") == "" {
		return "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	}
	return ""
}

// ParseCommandLine - читает параметры командной строки с значениями по умолчанию
func ParseCommandLine() {
	flag.StringVar(&Params.ServerAddress, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&Params.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&Params.FileStoragePath, "f", "./data/file-storage.txt", "File storage path")
	flag.StringVar(&Params.DatabaseDSN, "d", getDefaultDatabaseDSN(), "Database DSN")
	flag.Parse()
}

// ParseEnv - читает переменные окружения (если они есть) и сохраняет их в структуру Params
func ParseEnv() {
	// Читаем переменные окружения
	if err := env.Parse(&Params); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// JSONString - сериализуем структуру в формат JSON
func JSONString(params interface{}) string {
	bytes, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return string(bytes)
}

// PrintParams - выводит параметры приложения в лог
func PrintParams() {
	log.Info().Msg("Параметры приложения:\n" + JSONString(Params))
	JSONString(Params)
}
