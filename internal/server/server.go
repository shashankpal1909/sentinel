package server

import (
	"log/slog"
	"net/http"

	"sentinel/internal/app"
	"sentinel/internal/domain"
	"sentinel/internal/middleware"
	"sentinel/internal/proxy"
	"sentinel/internal/router"
)

type Server struct {
	router  *router.Router
	proxy   *proxy.Proxy
	handler http.Handler
	logger  *slog.Logger
}

func New(r *app.Runtime, p *proxy.Proxy, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	var routes []*domain.Route
	if r != nil {
		routes = r.Routes
	}
	if p == nil {
		p = proxy.New(logger)
	}
	s := &Server{
		router: router.New(routes),
		proxy:  p,
		logger: logger,
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
	service, ok := s.router.Match(r.URL.Path)
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
