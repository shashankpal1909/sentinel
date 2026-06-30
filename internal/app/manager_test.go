package app_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/domain"
)

type mockHealthUpdater struct {
	updateCount int
	lastRt      *app.Runtime
}

func (m *mockHealthUpdater) UpdateRuntime(ctx context.Context, newRt *app.Runtime) {
	m.updateCount++
	m.lastRt = newRt
}

func TestNewManager_Success(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://localhost:8001"},
				HealthCheck: &config.HealthCheckConfig{
					Path:               "/healthz",
					Interval:           "10s",
					Timeout:            "2s",
					HealthyThreshold:   1,
					UnhealthyThreshold: 2,
				},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/auth", Service: "auth"},
		},
	}

	mgr, err := app.NewManager(cfg, "gateway.yaml", slog.Default())
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if mgr.GetConfig() != cfg {
		t.Errorf("expected config to be stored")
	}
	rt := mgr.GetRuntime()
	if rt == nil || rt.Router == nil {
		t.Fatalf("expected runtime and router to be initialized")
	}
}

func TestManager_ApplyAndStatePreservation(t *testing.T) {
	initialCfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://localhost:8001", "http://localhost:8002"},
				HealthCheck: &config.HealthCheckConfig{
					Path:               "/healthz",
					Interval:           "10s",
					Timeout:            "2s",
					HealthyThreshold:   1,
					UnhealthyThreshold: 2,
				},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/auth", Service: "auth"},
		},
	}

	mgr, err := app.NewManager(initialCfg, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Manually set backend 8001 to unhealthy in old runtime
	rt := mgr.GetRuntime()
	rt.Services["auth"].Backends[0].SetState(domain.BackendStateUnhealthy)

	mockUpdater := &mockHealthUpdater{}
	mgr.SetHealthUpdater(context.Background(), mockUpdater)

	newCfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://localhost:8001", "http://localhost:8003"},
				HealthCheck: &config.HealthCheckConfig{
					Path:               "/healthz",
					Interval:           "10s",
					Timeout:            "2s",
					HealthyThreshold:   1,
					UnhealthyThreshold: 2,
				},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/auth", Service: "auth"},
		},
	}

	if err := mgr.Apply(newCfg); err != nil {
		t.Fatalf("unexpected apply error: %v", err)
	}

	if mockUpdater.updateCount != 1 {
		t.Errorf("expected health updater to be called once, got %d", mockUpdater.updateCount)
	}

	newRt := mgr.GetRuntime()
	if newRt == rt {
		t.Errorf("expected runtime pointer to be swapped")
	}

	// Verify state preservation: 8001 should be Unhealthy, 8003 should be Healthy (default)
	b1State := newRt.Services["auth"].Backends[0].GetState()
	if b1State != domain.BackendStateUnhealthy {
		t.Errorf("expected localhost:8001 state to be preserved as unhealthy, got %v", b1State)
	}
	b2State := newRt.Services["auth"].Backends[1].GetState()
	if b2State != domain.BackendStateHealthy {
		t.Errorf("expected localhost:8003 state to default to healthy, got %v", b2State)
	}
}

func TestManager_ReloadFromDisk(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "gateway.yaml")
	yamlContent := `
server:
  port: 8080
services:
  test:
    strategy: round-robin
    backends:
      - http://test:1234
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /test
    service: test
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	initialCfg, _ := config.Load(configPath)
	mgr, err := app.NewManager(initialCfg, configPath, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Update file on disk
	updatedYAML := `
server:
  port: 8080
services:
  test:
    strategy: round-robin
    backends:
      - http://test:1234
      - http://test:5678
    health_check:
      path: /healthz
      interval: 10s
      timeout: 2s
      healthy_threshold: 1
      unhealthy_threshold: 2
routes:
  - path: /test
    service: test
`
	if err := os.WriteFile(configPath, []byte(updatedYAML), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if err := mgr.ReloadFromDisk(); err != nil {
		t.Fatalf("unexpected reload error: %v", err)
	}

	if len(mgr.GetRuntime().Services["test"].Backends) != 2 {
		t.Errorf("expected 2 backends after reload, got %d", len(mgr.GetRuntime().Services["test"].Backends))
	}
}
