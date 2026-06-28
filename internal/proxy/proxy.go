package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"sync"

	"sentinel/internal/domain"
)

type Proxy struct {
	cache  sync.Map
	logger *slog.Logger
}

func New(logger *slog.Logger) *Proxy {
	if logger == nil {
		logger = slog.Default()
	}
	return &Proxy{
		logger: logger,
	}
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request, backend *domain.Backend) {
	if backend == nil || backend.URL == nil {
		p.logger.Error("Proxy forwarding failed: target backend is nil")
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	targetURL := backend.URL.String()

	val, ok := p.cache.Load(targetURL)
	if !ok {
		rp := httputil.NewSingleHostReverseProxy(backend.URL)
		rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			p.logger.Error("Upstream proxy request failed", "target", targetURL, "err", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
		val, _ = p.cache.LoadOrStore(targetURL, rp)
	}

	proxy := val.(*httputil.ReverseProxy)
	proxy.ServeHTTP(w, r)
}
