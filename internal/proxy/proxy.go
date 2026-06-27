package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"

	"sentinel/internal/domain"
)

type Proxy struct{}

func New() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Forward(backend *domain.Backend, w http.ResponseWriter, r *http.Request) {
	// Guard against nil backend targets to prevent dereference panics
	if backend == nil || backend.URL == nil {
		slog.Error("Proxy forwarding failed: target backend is nil")
		http.Error(w, "502 Bad Gateway: no upstream target available", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// Intercept transport errors to log details and return a clean 502 response
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("Upstream proxy request failed", "target", backend.URL.String(), "err", err)
		http.Error(w, "502 Bad Gateway: upstream server unreachable", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}
