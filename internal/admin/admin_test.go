package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"sentinel/internal/admin"
	"sentinel/internal/app"
	"sentinel/internal/config"
)

func setupTestServer(t *testing.T) (*admin.Server, *app.Manager, string) {
	t.Helper()
	dir := t.TempDir()
	configPath := filepath.Join(dir, "gateway.yaml")
	yamlContent := `
server:
  port: 8080
admin:
  port: 9901
services:
  auth:
    strategy: round-robin
    backends:
      - http://localhost:8001
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /auth
    service: auth
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cfg, _ := config.Load(configPath)
	mgr, err := app.NewManager(cfg, configPath, nil)
	if err != nil {
		t.Fatalf("unexpected manager error: %v", err)
	}

	srv := admin.New(mgr, nil)
	return srv, mgr, configPath
}

func TestAdmin_ConfigDump(t *testing.T) {
	srv, _, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/config_dump", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var dumped config.Config
	if err := json.Unmarshal(rec.Body.Bytes(), &dumped); err != nil {
		t.Fatalf("failed to unmarshal dump: %v", err)
	}
	if dumped.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got %d", dumped.Server.Port)
	}
}

func TestAdmin_ClustersAndListeners(t *testing.T) {
	srv, _, _ := setupTestServer(t)

	// Test /clusters
	req := httptest.NewRequest(http.MethodGet, "/clusters", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /clusters, got %d", rec.Code)
	}

	var clResp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &clResp)
	clusters := clResp["clusters"].([]interface{})
	if len(clusters) != 1 {
		t.Errorf("expected 1 cluster, got %d", len(clusters))
	}

	// Test /listeners
	req = httptest.NewRequest(http.MethodGet, "/listeners", nil)
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /listeners, got %d", rec.Code)
	}
}

func TestAdmin_ConfigApply(t *testing.T) {
	srv, mgr, _ := setupTestServer(t)

	newYAML := []byte(`
server:
  port: 8080
services:
  auth:
    strategy: round-robin
    backends:
      - http://localhost:8001
      - http://localhost:8002
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /auth
    service: auth
`)

	req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBuffer(newYAML))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	if len(mgr.GetRuntime().Services["auth"].Backends) != 2 {
		t.Errorf("expected 2 backends after dynamic apply, got %d", len(mgr.GetRuntime().Services["auth"].Backends))
	}
}

func TestAdmin_Reload(t *testing.T) {
	srv, mgr, configPath := setupTestServer(t)

	updatedYAML := `
server:
  port: 8080
services:
  auth:
    strategy: round-robin
    backends:
      - http://localhost:9000
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /auth
    service: auth
`
	os.WriteFile(configPath, []byte(updatedYAML), 0644)

	req := httptest.NewRequest(http.MethodPost, "/reload", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	if mgr.GetRuntime().Services["auth"].Backends[0].URL.String() != "http://localhost:9000" {
		t.Errorf("expected reloaded backend http://localhost:9000, got %s", mgr.GetRuntime().Services["auth"].Backends[0].URL.String())
	}
}
