package admin

import (
	"encoding/json"
	"net/http"

	"sentinel/internal/config"
	"sentinel/internal/domain"
)

func (s *Server) handleClusters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	snap := s.mgr.Current()
	var clusters []clusterResponse

	if snap != nil && snap.Runtime != nil && snap.Runtime.Services != nil {
		for name, svc := range snap.Runtime.Services {
			if svc == nil {
				continue
			}
			backends := make([]backendResponse, 0, len(svc.Backends))
			for _, b := range svc.Backends {
				if b != nil && b.URL != nil {
					backends = append(backends, backendResponse{
						URL:     b.URL.String(),
						State:   b.GetState().String(),
						Healthy: b.GetState() == domain.BackendStateHealthy,
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
			if cfg := s.mgr.CurrentConfig(); cfg != nil {
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
