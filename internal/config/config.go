// Инкремент 4
// Задание по треку «Сервис сокращения URL»
// Добавьте возможность конфигурировать сервис с помощью аргументов командной строки.
// Создайте конфигурацию или переменные для запуска со следующими флагами:
// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
// Флаг -b отвечает за базовый адрес результирующего сокращённого URL
// (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
// Совет: создайте отдельный пакет config, где будет храниться структура с вашей конфигурацией
// и функция, которая будет инициализировать поля этой структуры.
// По мере усложнения конфигурации вы сможете добавлять необходимые поля в вашу структуру и инициализировать их.

package config

import (
	"flag"
	"fmt"
)

var Address string
var BaseURL string

func ParseCommandLine() {
	flag.StringVar(&Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.Parse()

	fmt.Println("Server Address:", Address)
	fmt.Println("Shortened Base URL:", BaseURL)
}
