package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sentinel/internal/config"
)

func TestLoad_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_gateway.yaml")

	yamlContent := `
server:
  port: 9000
services:
  test-svc:
    strategy: round-robin
    backends:
      - http://localhost:8080
routes:
  - path: /test
    service: test-svc
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("expected successful load, got error: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("expected port 9000, got %d", cfg.Server.Port)
	}
	if len(cfg.Services) != 1 || len(cfg.Routes) != 1 {
		t.Errorf("expected 1 service and 1 route, got %d services and %d routes", len(cfg.Services), len(cfg.Routes))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("non-existent-file.yaml")
	if err == nil {
		t.Errorf("expected error loading non-existent file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	if err := os.WriteFile(configPath, []byte("invalid_yaml: [unclosed"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	_, err := config.Load(configPath)
	if err == nil {
		t.Errorf("expected error loading invalid YAML, got nil")
	}
}

func TestConfig_String(t *testing.T) {
	var nilCfg *config.Config
	if nilCfg.String() != "<nil>" {
		t.Errorf("expected <nil> string for nil config, got %s", nilCfg.String())
	}

	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
	}
	str := cfg.String()
	if !strings.Contains(str, "8080") {
		t.Errorf("expected string output to contain port 8080, got %s", str)
	}
}
