package middleware

import (
	"log"
	"net/http"
)

func Recover(logger *log.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					logger.Println("[ERROR] Panic caught:", r)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error\n"))
				}
			}()
			handler.ServeHTTP(w, req)
		})
	}
}
