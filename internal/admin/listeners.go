package admin

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleListeners(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	snap := s.mgr.Current()
	var listeners []listenerResponse
	if snap != nil && snap.Runtime != nil && snap.Runtime.Routes != nil {
		for _, route := range snap.Runtime.Routes {
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
