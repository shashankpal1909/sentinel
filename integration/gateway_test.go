package integration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/proxy"
	"sentinel/internal/router"
	"sentinel/internal/server"
)

func buildTestServer(t *testing.T, cfg *config.Config) *httptest.Server {
	t.Helper()
	rt, err := app.Build(cfg)
	if err != nil {
		t.Fatalf("failed to build runtime: %v", err)
	}
	snap := &app.Snapshot{Runtime: rt, Router: router.New(rt.Routes), Version: 1}
	srv := server.New(snap, proxy.New(nil), nil)
	return httptest.NewServer(srv)
}

func TestGateway_RouteResolution(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from route resolution"))
	}))
	defer backend.Close()

	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{
			"hello-svc": {Strategy: config.RoundRobin, Backends: []string{backend.URL}},
		},
		Routes: []config.RouteConfig{
			{Path: "/api/hello", Service: "hello-svc"},
		},
	}

	gw := buildTestServer(t, cfg)
	defer gw.Close()

	res, err := http.Get(gw.URL + "/api/hello")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "hello from route resolution" {
		t.Errorf("expected body 'hello from route resolution', got %q", string(body))
	}
}

func TestGateway_RoundRobin(t *testing.T) {
	b1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("node-1"))
	}))
	defer b1.Close()

	b2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("node-2"))
	}))
	defer b2.Close()

	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{
			"rr-svc": {Strategy: config.RoundRobin, Backends: []string{b1.URL, b2.URL}},
		},
		Routes: []config.RouteConfig{
			{Path: "/rr", Service: "rr-svc"},
		},
	}

	gw := buildTestServer(t, cfg)
	defer gw.Close()

	expected := []string{"node-1", "node-2", "node-1", "node-2"}
	for i, exp := range expected {
		res, err := http.Get(gw.URL + "/rr")
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		if string(body) != exp {
			t.Errorf("request %d: expected %q, got %q", i, exp, string(body))
		}
	}
}

func TestGateway_RandomDistribution(t *testing.T) {
	b1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("rnd-1")) }))
	defer b1.Close()
	b2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("rnd-2")) }))
	defer b2.Close()
	b3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("rnd-3")) }))
	defer b3.Close()

	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{
			"rnd-svc": {Strategy: config.Random, Backends: []string{b1.URL, b2.URL, b3.URL}},
		},
		Routes: []config.RouteConfig{
			{Path: "/rnd", Service: "rnd-svc"},
		},
	}

	gw := buildTestServer(t, cfg)
	defer gw.Close()

	counts := make(map[string]int)
	for i := 0; i < 100; i++ {
		res, err := http.Get(gw.URL + "/rnd")
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		counts[string(body)]++
	}

	for _, expectedNode := range []string{"rnd-1", "rnd-2", "rnd-3"} {
		if counts[expectedNode] == 0 {
			t.Errorf("expected node %q to receive at least 1 request in 100 random hits, got 0", expectedNode)
		}
	}
}

func TestGateway_NotFound(t *testing.T) {
	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{},
		Routes:   []config.RouteConfig{},
	}

	gw := buildTestServer(t, cfg)
	defer gw.Close()

	res, err := http.Get(gw.URL + "/non-existent-path")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found, got %d", res.StatusCode)
	}
}

func TestGateway_BadGateway(t *testing.T) {
	// Create backend and immediately shut it down so connection attempts fail
	deadBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	deadURL := deadBackend.URL
	deadBackend.Close()

	cfg := &config.Config{
		Services: map[string]config.ServiceConfig{
			"dead-svc": {Strategy: config.RoundRobin, Backends: []string{deadURL}},
		},
		Routes: []config.RouteConfig{
			{Path: "/fail", Service: "dead-svc"},
		},
	}

	gw := buildTestServer(t, cfg)
	defer gw.Close()

	res, err := http.Get(gw.URL + "/fail")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway when upstream fails, got %d", res.StatusCode)
	}
}
