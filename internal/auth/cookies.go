package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/shortener"
)

// contextKey - тип для ключа контекста
type contextKey string

const UserIDKey contextKey = "userID"
const NewUserIDKey contextKey = "newUserID"

// GenerateUserID - генерирует и возвращает новый случайный ID пользователя
func GenerateUserID() string {
	// return Params.UserID
	return "us-" + shortener.Shorten(uuid.New().String())
}

// AuthMiddleware - middleware для аутентификации пользователя
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получить куки из запроса
		cookie, err := r.Cookie(Params.CookieName)
		// Если куки отсутствует, создаём новые и добавляем в ответ перед продолжением обработки запроса
		if err != nil {
			log.Warn().Msg("AuthMiddleware> Cookie '%v' not found.  Adding New user ID to response and request")
			ctx := addCookieAndContext(w, r)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Получаем userID и подпись из куки
		userID, signature := getUserIDAndSignature(cookie.Value)

		// Если подпись неверна, создаём новые куки и добавляем в ответ перед продолжением обработки запроса
		if !signatureValid(userID, signature) {
			log.Warn().Msgf("AuthMiddleware> Invalid cookie signature for user ID '%v'. Adding a new cookie '%v' to the response and request", userID, Params.CookieName)
			ctx := addCookieAndContext(w, r)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Добавляем ID пользователя в контекст запроса
		ctx := AddUserIDToContext(r.Context(), userID, "old")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Функция для создания подписанной строки куки
func signCookie(userID string) string {
	signature := computeHMAC(userID, []byte(Params.SecretKey))
	return userID + "|" + signature
}

// Функция для вычисления HMAC подписи
func computeHMAC(message string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// SetResponseCookie - устанавливает куки в ответе
// Подписывает значение куки и добавляет в ответ.
// Параметры:
//   - w - http.ResponseWriter
//   - cookieValue - значение куки
func SetResponseCookie(w http.ResponseWriter, cookieValue string) {
	signedCookie := signCookie(cookieValue)
	http.SetCookie(w, &http.Cookie{
		Name:  Params.CookieName,
		Value: signedCookie,
		Path:  "/",
		// Установите дополнительные параметры безопасности по необходимости
		// HttpOnly: true,
		// Secure:   true,
	})
	log.Info().Msgf(">>> SetResponseCookie> Setting cookie '%v' with value '%v'", Params.CookieName, signedCookie[:20]+"...")
}

// AddUserIDToContext - Добавляем ID пользователя, и метку новый ли он в контекст запроса
// Параметры:
//   - ctx - контекст
//   - newUserID - новый ID пользователя
//   - keyLabel - метка нового ID
func AddUserIDToContext(ctx context.Context, newUserID, keyLabel string) (newCtx context.Context) {
	newCtx = context.WithValue(ctx, UserIDKey, newUserID)
	newCtx = context.WithValue(newCtx, NewUserIDKey, keyLabel)
	log.Info().Msgf(">>> AddUserIDToContext> New User ID '%v' is added to request context", newUserID)
	return newCtx
}

// getUserIDAndSignature - Разделяем значение куки на ID и подпись
// Параметры:
//   - cookieValue - значение куки
//
// Возвращает:
//   - userID - ID пользователя
//   - signature - подпись
func getUserIDAndSignature(cookieValue string) (userID, signature string) {
	parts := strings.Split(cookieValue, "|")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// signatureValid - Проверяет UserID и подпись
// Параметры:
//   - userID - ID пользователя
//   - signature - подпись
//
// Возвращает:
//   - true, если подпись верна
func signatureValid(userID, signature string) bool {
	if userID == "" {
		return false
	}
	expectedSignature := computeHMAC(userID, []byte(Params.SecretKey))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// addCookieAndContext - Генерирует новый ID пользователя,
// устанавливает куки в ответе сервера и добавляет ID пользователя в контекст запроса.
// Параметры:
//   - w - http.ResponseWriter
//   - r - *http.Request
//
// Возвращает:
//   - контекст запроса
func addCookieAndContext(w http.ResponseWriter, r *http.Request) context.Context {
	// Генерируем новый ID пользователя
	newUserID := GenerateUserID()
	// Устанавливаем куки в ответе сервера
	SetResponseCookie(w, newUserID)
	// Добавляем ID пользователя в контекст запроса
	ctx := AddUserIDToContext(r.Context(), newUserID, "new")
	return ctx
}
