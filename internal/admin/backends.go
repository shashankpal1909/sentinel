package admin

import (
	"encoding/json"
	"net/http"
)

type backendDetailResponse struct {
	Service            string `json:"service"`
	URL                string `json:"url"`
	State              string `json:"state"`
	Interval           string `json:"interval"`
	HealthyThreshold   int    `json:"healthyThreshold"`
	UnhealthyThreshold int    `json:"unhealthyThreshold"`
}

func (s *Server) handleBackends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	snap := s.mgr.Current()
	backends := make([]backendDetailResponse, 0)

	if snap != nil && snap.Runtime != nil && snap.Runtime.Services != nil {
		for name, svc := range snap.Runtime.Services {
			if svc == nil {
				continue
			}
			interval := "N/A"
			healthyThresh := 0
			unhealthyThresh := 0
			if svc.Health.Path != "" {
				interval = svc.Health.Interval.String()
				healthyThresh = svc.Health.HealthyThreshold
				unhealthyThresh = svc.Health.UnhealthyThreshold
			}
			for _, b := range svc.Backends {
				if b != nil && b.URL != nil {
					backends = append(backends, backendDetailResponse{
						Service:            name,
						URL:                b.URL.String(),
						State:              b.GetState().String(),
						Interval:           interval,
						HealthyThreshold:   healthyThresh,
						UnhealthyThreshold: unhealthyThresh,
					})
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"backends": backends,
	})
}
