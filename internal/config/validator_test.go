package config_test

import (
	"testing"

	"sentinel/internal/config"
)

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://localhost:8001"},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/auth", Service: "auth"},
		},
	}

	if err := config.Validate(cfg); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 0},
	}
	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected error for port 0, got nil")
	}

	cfg.Server.Port = 70000
	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected error for port 70000, got nil")
	}
}

func TestValidate_EmptyBackends(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{},
			},
		},
	}
	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected error for empty backends, got nil")
	}
}

func TestValidate_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"invalid-url"},
			},
		},
	}
	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected error for invalid backend URL, got nil")
	}
}

func TestValidate_DuplicateRoutes(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://localhost:8001"},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/auth", Service: "auth"},
			{Path: "/auth", Service: "auth"},
		},
	}
	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected error for duplicate routes, got nil")
	}
}

func TestValidate_MissingReference(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://localhost:8001"},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/users", Service: "user-service"},
		},
	}
	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected error for missing service reference, got nil")
	}
}
