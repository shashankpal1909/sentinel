package server

import (
	"log/slog"
	"net/http"

	"sentinel/internal/app"
	"sentinel/internal/domain"
	"sentinel/internal/proxy"
	"sentinel/internal/router"
)

type Server struct {
	router *router.Router
	proxy  *proxy.Proxy
}

func New(r *app.Runtime, p *proxy.Proxy) *Server {
	var routes []*domain.Route
	if r != nil {
		routes = r.Routes
	}
	if p == nil {
		p = proxy.New()
	}
	return &Server{
		router: router.New(routes),
		proxy:  p,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s == nil || s.router == nil || s.proxy == nil {
		http.Error(w, "500 Internal Server Error: server uninitialized", http.StatusInternalServerError)
		return
	}

	service, ok := s.router.Match(r.URL.Path)
	if !ok || service == nil {
		slog.Warn("No service found for path", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	backend, err := service.NextBackend()
	if err != nil {
		slog.Error("Error getting backend", "service", service.Name, "err", err)
		http.Error(w, "502 Bad Gateway: no upstream target available", http.StatusBadGateway)
		return
	}

	s.proxy.Forward(w, r, backend)
}
