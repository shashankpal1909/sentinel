package middleware

import "net/http"

func Chain(
	h http.Handler,
	middlewares ...Middleware,
) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}
