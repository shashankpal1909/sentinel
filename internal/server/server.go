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
}

func New(r *app.Runtime, p *proxy.Proxy) *Server {
	var routes []*domain.Route
	if r != nil {
		routes = r.Routes
	}
	if p == nil {
		p = proxy.New()
	}
	s := &Server{
		router: router.New(routes),
		proxy:  p,
	}
	s.handler = middleware.Chain(
		http.HandlerFunc(s.handleRoute),
		middleware.Recovery(),
		middleware.Logger(),
		middleware.RequestID(),
	)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s == nil || s.handler == nil {
		http.Error(w, "500 Internal Server Error: server uninitialized", http.StatusInternalServerError)
		return
	}
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleRoute(w http.ResponseWriter, r *http.Request) {
	if s.router == nil || s.proxy == nil {
		http.Error(w, "500 Internal Server Error: server uninitialized", http.StatusInternalServerError)
		return
	}

	reqID := middleware.GetRequestID(r.Context())
	service, ok := s.router.Match(r.URL.Path)
	if !ok || service == nil {
		slog.Warn("No service found for path", "path", r.URL.Path, "request_id", reqID)
		http.NotFound(w, r)
		return
	}

	backend, err := service.NextBackend()
	if err != nil {
		slog.Error("Error getting backend", "service", service.Name, "request_id", reqID, "err", err)
		http.Error(w, "502 Bad Gateway: no upstream target available", http.StatusBadGateway)
		return
	}

	s.proxy.Forward(w, r, backend)
}
