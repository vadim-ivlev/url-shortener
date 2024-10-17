// Description: функции для работы с JWT-токенами.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildToken создаёт токен и возвращает его в виде строки.
func BuildToken(userID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Params.TokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(Params.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// keyFunc — функция для получения ключа подписи
func keyFunc(t *jwt.Token) (interface{}, error) {
	// проверяем, что используется алгоритм HS256
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	return []byte(Params.SecretKey), nil
}

// GetUserIDFromToken возвращает ID пользователя из токена JWT
func GetUserIDFromToken(jwtTokenString string) (userID string, err error) {
	// создаём экземпляр структуры с утверждениями
	claims := &Claims{}

	// парсим из строки токена jwtTokenString в структуру claims
	token, err := jwt.ParseWithClaims(jwtTokenString, claims, keyFunc)
	if err != nil {
		return "", err
	}

	// проверяем, что токен валиден
	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	return claims.UserID, nil
}
