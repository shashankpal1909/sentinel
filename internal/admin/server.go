package admin

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"

	"sentinel/internal/app"
)

type Server struct {
	mgr        *app.Manager
	loader     *app.Loader
	configPath string
	health     HealthUpdater
	healthCtx  context.Context
	logger     *slog.Logger
	mux        *http.ServeMux
	version    atomic.Uint64
}

// Option configures Server.
type Option func(*Server)

func WithLoader(loader *app.Loader) Option {
	return func(s *Server) { s.loader = loader }
}

func WithConfigPath(path string) Option {
	return func(s *Server) { s.configPath = path }
}

func New(mgr *app.Manager, logger *slog.Logger, opts ...Option) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{
		mgr:    mgr,
		loader: app.NewLoader(),
		logger: logger,
		mux:    http.NewServeMux(),
	}
	if mgr != nil {
		snap := mgr.Current()
		if snap != nil {
			s.version.Store(snap.Version)
		}
	}
	for _, opt := range opts {
		opt(s)
	}
	s.registerRoutes()
	return s
}

func (s *Server) SetHealthUpdater(ctx context.Context, h HealthUpdater) {
	s.healthCtx = ctx
	s.health = h
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/config_dump", s.handleConfigDump)
	s.mux.HandleFunc("/clusters", s.handleClusters)
	s.mux.HandleFunc("/services", s.handleClusters) // alias
	s.mux.HandleFunc("/listeners", s.handleListeners)
	s.mux.HandleFunc("/routes", s.handleListeners) // alias
	s.mux.HandleFunc("/config", s.handleConfigApply)
	s.mux.HandleFunc("/reload", s.handleReload)
	s.mux.HandleFunc("/runtime", s.handleRuntime)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleConfigDump(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := s.mgr.CurrentConfig()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (s *Server) handleRuntime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	snap := s.mgr.Current()
	resp := runtimeResponse{
		Version:  0,
		LoadedAt: "",
		Services: 0,
		Routes:   0,
	}
	if snap != nil {
		resp.Version = snap.Version
		if !snap.LoadedAt.IsZero() {
			resp.LoadedAt = snap.LoadedAt.Format("2006-01-02T15:04:05Z")
		}
		if snap.Runtime != nil {
			resp.Services = len(snap.Runtime.Services)
			resp.Routes = len(snap.Runtime.Routes)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
