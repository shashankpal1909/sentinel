package admin

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"go.yaml.in/yaml/v4"
	"sentinel/internal/app"
	"sentinel/internal/config"
)

type Server struct {
	mgr    *app.Manager
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(mgr *app.Manager, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{
		mgr:    mgr,
		logger: logger,
		mux:    http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/config_dump", s.handleConfigDump)
	s.mux.HandleFunc("/clusters", s.handleClusters)
	s.mux.HandleFunc("/services", s.handleClusters) // alias
	s.mux.HandleFunc("/listeners", s.handleListeners)
	s.mux.HandleFunc("/routes", s.handleListeners) // alias
	s.mux.HandleFunc("/config", s.handleConfigApply)
	s.mux.HandleFunc("/reload", s.handleReload)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleConfigDump(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := s.mgr.GetConfig()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

type backendResponse struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

type clusterResponse struct {
	Name        string                    `json:"name"`
	Strategy    config.BalancerStrategy   `json:"strategy"`
	HealthCheck *config.HealthCheckConfig `json:"health_check,omitempty"`
	Backends    []backendResponse         `json:"backends"`
}

func (s *Server) handleClusters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	rt := s.mgr.GetRuntime()
	var clusters []clusterResponse

	if rt != nil && rt.Services != nil {
		for name, svc := range rt.Services {
			if svc == nil {
				continue
			}
			backends := make([]backendResponse, 0, len(svc.Backends))
			for _, b := range svc.Backends {
				if b != nil && b.URL != nil {
					backends = append(backends, backendResponse{
						URL:   b.URL.String(),
						State: b.GetState().String(),
					})
				}
			}
			var hcConfig *config.HealthCheckConfig
			if svc.Health.Path != "" {
				hcConfig = &config.HealthCheckConfig{
					Path:               svc.Health.Path,
					Interval:           svc.Health.Interval.String(),
					Timeout:            svc.Health.Timeout.String(),
					HealthyThreshold:   svc.Health.HealthyThreshold,
					UnhealthyThreshold: svc.Health.UnhealthyThreshold,
				}
			}
			strategy := config.BalancerStrategy("unknown")
			if cfg := s.mgr.GetConfig(); cfg != nil {
				if svcCfg, ok := cfg.Services[name]; ok && svcCfg.Strategy != "" {
					strategy = svcCfg.Strategy
				}
			}
			clusters = append(clusters, clusterResponse{
				Name:        name,
				Strategy:    strategy,
				HealthCheck: hcConfig,
				Backends:    backends,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"clusters": clusters,
	})
}

type listenerResponse struct {
	Path    string `json:"path"`
	Service string `json:"service"`
}

func (s *Server) handleListeners(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	rt := s.mgr.GetRuntime()
	var listeners []listenerResponse
	if rt != nil && rt.Routes != nil {
		for _, route := range rt.Routes {
			if route != nil {
				svcName := ""
				if route.Service != nil {
					svcName = route.Service.Name
				}
				listeners = append(listeners, listenerResponse{
					Path:    route.Path,
					Service: svcName,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"listeners": listeners,
	})
}

func (s *Server) handleConfigApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newCfg config.Config
	if err := yaml.Unmarshal(body, &newCfg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse configuration payload: " + err.Error(),
		})
		return
	}

	if err := s.mgr.Apply(&newCfg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "applied",
		"message": "Configuration successfully applied via hot reload",
	})
}

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.mgr.ReloadFromDisk(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "reloaded",
		"message": "Configuration successfully reloaded from disk",
	})
}
