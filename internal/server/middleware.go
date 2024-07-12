package server

import "net/http"

// middleware для установки Content-Type в значение application/json
// https://github.com/oapi-codegen/oapi-codegen/issues/97
func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
