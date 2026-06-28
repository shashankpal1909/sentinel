package middleware

import (
	"log/slog"
	"net/http"
)

func Recovery(logger *slog.Logger) Middleware {
	if logger == nil {
		logger = slog.Default()
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					logger.Error("panic", "err", rcv)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
