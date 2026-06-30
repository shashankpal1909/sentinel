package server

import (
	"log/slog"
	"net/http"

	"sentinel/internal/app"
	"sentinel/internal/middleware"
	"sentinel/internal/proxy"
)

type SnapshotProvider interface {
	Current() *app.Snapshot
}

type Server struct {
	provider SnapshotProvider
	proxy    *proxy.Proxy
	handler  http.Handler
	logger   *slog.Logger
}

func New(provider SnapshotProvider, p *proxy.Proxy, logger *slog.Logger) *Server {
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
	var snap *app.Snapshot
	if s.provider != nil {
		snap = s.provider.Current()
	}
	if snap == nil || snap.Router == nil {
		s.logger.Warn("No snapshot found", "path", r.URL.Path, "request_id", reqID)
		http.NotFound(w, r)
		return
	}

	service, ok := snap.Router.Match(r.URL.Path)
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
