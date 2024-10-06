// Description: Этот пакет содержит функции для аутентификации пользователей.
package auth

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type authParams struct {
	// время жизни JWT токена
	TokenExp time.Duration `env:"TOKEN_EXP" envDefault:"3h"`
	// секретный ключ JWT токена
	SecretKey string `env:"SECRET_KEY " envDefault:"supersecret"`
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
