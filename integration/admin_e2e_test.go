package integration_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"sentinel/internal/admin"
	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/proxy"
	"sentinel/internal/server"
)

func TestGatewayAdmin_HotReloadE2E(t *testing.T) {
	b1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("version-1"))
	}))
	defer b1.Close()

	b2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("version-2"))
	}))
	defer b2.Close()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "gateway.yaml")
	initialYAML := `
server:
  port: 8080
admin:
  port: 9901
services:
  app-svc:
    strategy: round-robin
    backends:
      - ` + b1.URL + `
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /app
    service: app-svc
`
	if err := os.WriteFile(configPath, []byte(initialYAML), 0644); err != nil {
		t.Fatalf("failed to write initial config: %v", err)
	}

	cfg, _ := config.Load(configPath)
	mgr, err := app.NewManager(cfg, configPath, nil)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	gwServer := server.New(mgr, proxy.New(nil), nil)
	gwTS := httptest.NewServer(gwServer)
	defer gwTS.Close()

	adminServer := admin.New(mgr, nil)
	adminTS := httptest.NewServer(adminServer)
	defer adminTS.Close()

	// 1. Verify Gateway routes to b1 ("version-1")
	res, err := http.Get(gwTS.URL + "/app")
	if err != nil {
		t.Fatalf("gateway request 1 failed: %v", err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if string(body) != "version-1" {
		t.Fatalf("expected 'version-1', got %q", string(body))
	}

	// 2. Push new config via Admin API POST /config pointing to b2 ("version-2")
	updatedYAML := []byte(`
server:
  port: 8080
admin:
  port: 9901
services:
  app-svc:
    strategy: round-robin
    backends:
      - ` + b2.URL + `
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /app
    service: app-svc
`)
	resAdmin, err := http.Post(adminTS.URL+"/config", "application/x-yaml", bytes.NewBuffer(updatedYAML))
	if err != nil {
		t.Fatalf("admin POST /config failed: %v", err)
	}
	resAdmin.Body.Close()
	if resAdmin.StatusCode != http.StatusOK {
		t.Fatalf("expected admin status 200 OK, got %d", resAdmin.StatusCode)
	}

	// 3. Verify Gateway now routes to b2 ("version-2") dynamically without restart
	res2, err := http.Get(gwTS.URL + "/app")
	if err != nil {
		t.Fatalf("gateway request 2 failed: %v", err)
	}
	body2, _ := io.ReadAll(res2.Body)
	res2.Body.Close()
	if string(body2) != "version-2" {
		t.Fatalf("expected 'version-2' after dynamic reload, got %q", string(body2))
	}
}
