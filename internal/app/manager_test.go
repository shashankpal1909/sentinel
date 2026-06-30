package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/domain"
)

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

	loader := app.NewLoader()
	snap, err := loader.Build(cfg, 1)
	if err != nil {
		t.Fatalf("expected build success, got error: %v", err)
	}

	mgr, err := app.NewManager(snap, cfg)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if mgr.GetConfig() != cfg {
		t.Errorf("expected config to be stored")
	}
	s := mgr.Current()
	if s == nil || s.Router == nil {
		t.Fatalf("expected snapshot and router to be initialized")
	}
}

func TestManager_ReplaceAndStatePreservation(t *testing.T) {
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

	loader := app.NewLoader()
	snap1, err := loader.Build(initialCfg, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mgr, err := app.NewManager(snap1, initialCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Manually set backend 8001 to unhealthy in old runtime
	rt := mgr.GetRuntime()
	rt.Services["auth"].Backends[0].SetState(domain.BackendStateUnhealthy)

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

	snap2, err := loader.Build(newCfg, 2)
	if err != nil {
		t.Fatalf("unexpected build error: %v", err)
	}

	mgr.Replace(snap2, newCfg)

	newSnap := mgr.Current()
	if newSnap == snap1 {
		t.Errorf("expected snapshot pointer to be swapped")
	}
	if newSnap.Version != 2 {
		t.Errorf("expected snapshot version 2, got %d", newSnap.Version)
	}

	newRt := mgr.GetRuntime()
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

func TestLoader_LoadFromDisk(t *testing.T) {
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

	loader := app.NewLoader()
	cfg, snap, err := loader.Load(configPath, 1)
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
	if len(snap.Runtime.Services["test"].Backends) != 1 {
		t.Errorf("expected 1 backend, got %d", len(snap.Runtime.Services["test"].Backends))
	}
}
