package server

import (
	"log/slog"
	"net/http"

	"sentinel/internal/app"
	"sentinel/internal/middleware"
	"sentinel/internal/proxy"
	"sentinel/internal/router"
)

type RuntimeProvider interface {
	GetRuntime() *app.Runtime
}

type Server struct {
	provider RuntimeProvider
	proxy    *proxy.Proxy
	handler  http.Handler
	logger   *slog.Logger
}

func New(provider RuntimeProvider, p *proxy.Proxy, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	if p == nil {
		p = proxy.New(logger)
	}
	s := &Server{
		provider: provider,
		proxy:    p,
		logger:   logger,
	}
	s.handler = middleware.Chain(
		http.HandlerFunc(s.handleRoute),
		middleware.Recovery(logger),
		middleware.RequestID(),
		middleware.Logger(logger),
	)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleRoute(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestID(r.Context())
	var rt *app.Runtime
	if s.provider != nil {
		rt = s.provider.GetRuntime()
	}
	if rt == nil {
		s.logger.Warn("No runtime found", "path", r.URL.Path, "request_id", reqID)
		http.NotFound(w, r)
		return
	}

	rtr := rt.Router
	if rtr == nil {
		rtr = router.New(rt.Routes)
	}

	service, ok := rtr.Match(r.URL.Path)
	if !ok || service == nil {
		s.logger.Warn("No service found for path", "path", r.URL.Path, "request_id", reqID)
		http.NotFound(w, r)
		return
	}

	backend, err := service.NextBackend()
	if err != nil {
		s.logger.Error("Error getting backend", "service", service.Name, "request_id", reqID, "err", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	s.proxy.Forward(w, r, backend)
}
