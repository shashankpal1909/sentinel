package middleware

import (
	"log/slog"
	"net/http"
)

func Recovery() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					slog.Error("panic", "err", rcv)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
