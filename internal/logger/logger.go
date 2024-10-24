package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Расширяем стандартный http.ResponseWriter для захвата данных об ответе
type loggingResponseWriter struct {
	// встраиваем оригинальный http.ResponseWriter
	http.ResponseWriter
	// код статуса ответа
	status int
	// размер ответа
	size int
	// content-type ответа
	contentType string
	// content-encoding ответа
	contentEncoding string
}

// переопределяем метод Write для захвата размера ответа
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, в оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	// захватываем размер ответа
	r.size += size
	// захватываем content-type ответа
	r.contentType = r.Header().Get("Content-Type")
	// захватываем content-encoding ответа
	r.contentEncoding = r.Header().Get("Content-Encoding")
	return size, err
}

// переопределяем метод WriteHeader для захвата кода статуса
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, в оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	// захватываем код статуса
	r.status = statusCode
}

// NoColor — флаг для отключения цветного лога
var NoColor bool = false

// инициализируем логгер
func InitializeLogger() {
	// цветной лог в консоль
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: NoColor})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Logger initialized")
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
// - Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
// - Сведения об ответах должны содержать код статуса и размер содержимого ответа.
func RequestLogger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// встраиваем оригинальный http.ResponseWriter в loggingResponseWriter
		logRespWriter := loggingResponseWriter{ResponseWriter: w}

		start := time.Now()

		// выполняем запрос с записью ответа в logRespWriter
		h.ServeHTTP(&logRespWriter, r)

		uri := r.RequestURI
		method := r.Method
		acceptEncoding := r.Header.Get("Accept-Encoding")
		duration := time.Since(start)
		log.Info().
			Str("URI", uri).
			Str("method", method).
			Dur("duration", duration).
			Str("accept-encoding", acceptEncoding).
			Int("size", logRespWriter.size).
			Int("status", logRespWriter.status).
			Str("content-type", logRespWriter.contentType).
			Str("content-encoding", logRespWriter.contentEncoding).
			Msg("")

	}

	return http.HandlerFunc(logFn)
}
