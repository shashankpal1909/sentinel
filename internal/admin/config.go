package admin

import (
	"encoding/json"
	"io"
	"net/http"

	"go.yaml.in/yaml/v4"
	"sentinel/internal/config"
)

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

	nextVersion := s.version.Add(1)
	snap, err := s.loader.Build(&newCfg, nextVersion)
	if err != nil {
		s.version.Add(^uint64(0)) // rollback version increment on error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.mgr.Replace(snap, &newCfg)
	if s.health != nil && s.healthCtx != nil {
		s.health.UpdateRuntime(s.healthCtx, snap.Runtime)
	}

	s.logger.Info("Configuration applied successfully via hot reload", "version", nextVersion)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "applied",
		"message": "Configuration successfully applied via hot reload",
	})
}
