package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"sync"

	"sentinel/internal/domain"
)

type Proxy struct {
	cache sync.Map
}

func New() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, backend *domain.Backend) {
	if backend == nil || backend.URL == nil {
		slog.Error("Proxy forwarding failed: target backend is nil")
		http.Error(w, "502 Bad Gateway: no upstream target available", http.StatusBadGateway)
		return
	}

	targetURL := backend.URL.String()

	// Look up cached ReverseProxy to avoid re-allocating immutable Director and Transport
	val, ok := p.cache.Load(targetURL)
	if !ok {
		rp := httputil.NewSingleHostReverseProxy(backend.URL)
		rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("Upstream proxy request failed", "target", targetURL, "err", err)
			http.Error(w, "502 Bad Gateway: upstream server unreachable", http.StatusBadGateway)
		}
		val, _ = p.cache.LoadOrStore(targetURL, rp)
	}

	proxy := val.(*httputil.ReverseProxy)
	proxy.ServeHTTP(w, r)
}
