package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// AuthMiddleware - middleware для аутентификации пользователя
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieName := Params.CookieName
		// Попытка получить куки из запроса
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			// Куки отсутствует, создаём новую и добавляем в ответ перед продолжением обработки запроса
			log.Warn().Msgf("Cookie '%v' not found. Adding a new one to the response", cookieName)
			newID := Params.UserID //uuid.New().String()
			log.Info().Msgf("New user ID is '%v' ", newID)
			signedCookie := signCookie(newID)
			http.SetCookie(w, &http.Cookie{
				Name:  cookieName,
				Value: signedCookie,
				Path:  "/",
				// Установите дополнительные параметры безопасности по необходимости
				// HttpOnly: true,
				// Secure:   true,
			})
			// Продолжаем обработку запроса
			next.ServeHTTP(w, r)
			return
		}

		// Разделяем значение куки на ID и подпись
		parts := strings.Split(cookie.Value, "|")
		if len(parts) != 2 {
			http.Error(w, "Invalid cookie format", http.StatusUnauthorized)
			return
		}
		userID := parts[0]
		signature := parts[1]

		// Проверяем подпись
		log.Info().Msgf("Checking cookie signature for user ID %v", userID)
		expectedSignature := computeHMAC(userID, []byte(Params.SecretKey))
		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			// Подпись неверна, создаём новую куку
			log.Warn().Msgf("Invalid cookie signature for user ID '%v'. Adding a new cookie '%v' to the response", userID, cookieName)
			newID := Params.UserID //uuid.New().String()
			signedCookie := signCookie(newID)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    signedCookie,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
			})
			// Продолжаем обработку запроса
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем, что ID существует
		if userID == "" {
			http.Error(w, "Unauthorized: No user ID", http.StatusUnauthorized)
			log.Error().Msg("No user ID in the cookie")
			return
		}

		// Можно добавить ID пользователя в контекст запроса, если необходимо
		ctx := context.WithValue(r.Context(), "userID", userID)
		log.Info().Msgf("User ID '%v' is authenticated and added to request context", userID)
		next.ServeHTTP(w, r.WithContext(ctx))

		// Продолжаем обработку запроса
		// next.ServeHTTP(w, r)
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

// // Пример обработчика, защищённого middleware
// func protectedHandler(w http.ResponseWriter, r *http.Request) {
// 	// Здесь можно получить userID из контекста, если вы его добавили
// 	// userID := r.Context().Value("userID").(string)
// 	w.Write([]byte("Доступ разрешён"))
// }

// func main() {
// 	mux := http.NewServeMux()
// 	// Применяем middleware к защищённому маршруту
// 	mux.Handle("/protected", authMiddleware(http.HandlerFunc(protectedHandler)))

// 	// Запуск сервера
// 	log.Println("Сервер запущен на :8080")
// 	if err := http.ListenAndServe(":8080", mux); err != nil {
// 		log.Fatalf("Ошибка запуска сервера: %v", err)
// 	}
// }
