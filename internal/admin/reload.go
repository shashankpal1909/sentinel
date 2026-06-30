package admin

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.configPath == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": errors.New("no configuration file path set for reload").Error(),
		})
		return
	}

	nextVersion := s.version.Add(1)
	newCfg, snap, err := s.loader.Load(s.configPath, nextVersion)
	if err != nil {
		s.version.Add(^uint64(0))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	s.mgr.Replace(snap, newCfg)
	if s.health != nil && s.healthCtx != nil {
		s.health.UpdateRuntime(s.healthCtx, snap.Runtime)
	}

	s.logger.Info("Configuration reloaded successfully from disk", "path", s.configPath, "version", nextVersion)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "reloaded",
		"message": "Configuration successfully reloaded from disk",
	})
}
