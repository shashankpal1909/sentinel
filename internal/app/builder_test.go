package app_test

import (
	"bytes"
	"testing"

	"sentinel/internal/app"
	"sentinel/internal/config"
)

func TestBuild_SharedServiceRegistry(t *testing.T) {
	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{
			"userService": {
				Strategy: "round-robin",
				Backends: []string{"http://user1:8080", "http://user2:8080"},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/users", Service: "userService"},
			{Path: "/profile", Service: "userService"},
		},
	}

	rt, err := app.Build(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(rt.Routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(rt.Routes))
	}

	if len(rt.Services) != 1 {
		t.Fatalf("expected 1 service in registry, got %d", len(rt.Services))
	}

	// Verify Shared Service Registry: both routes must point to the exact same Service pointer
	if rt.Routes[0].Service != rt.Routes[1].Service {
		t.Errorf("expected both routes to share the exact same Service instance pointer")
	}

	svc := rt.Services["userService"]
	if rt.Routes[0].Service != svc {
		t.Errorf("expected route service pointer to match Services map entry")
	}
}

func TestRuntime_DumpAndString(t *testing.T) {
	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{
			"auth": {
				Strategy: "round-robin",
				Backends: []string{"http://auth:8080"},
			},
		},
		Routes: []config.RouteConfig{
			{Path: "/login", Service: "auth"},
		},
	}

	rt, err := app.Build(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	rt.Dump(&buf)
	if buf.Len() == 0 {
		t.Errorf("expected Dump output, got empty buffer")
	}

	str := rt.String()
	if len(str) == 0 {
		t.Errorf("expected String output, got empty string")
	}
}
