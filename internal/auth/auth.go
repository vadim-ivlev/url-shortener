// Description: Этот пакет содержит функции для аутентификации пользователей.
package auth

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type authParams struct {
	// секретный ключ для подписи и проверки куки
	SecretKey string `env:"SECRET_KEY " envDefault:"supersecret"`
	// Имя куки
	CookieName string `env:"COOKIE_NAME" envDefault:"url-shortener"`
}

// Params - переменная для хранения параметров auth
var Params authParams = authParams{}

// init - инициализируем модуль
func init() {
	ParseEnv()
}

// ParseEnv - читает переменные окружения (если они есть) и сохраняет их в структуру Params
func ParseEnv() {
	// Читаем переменные окружения
	if err := env.Parse(&Params); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
