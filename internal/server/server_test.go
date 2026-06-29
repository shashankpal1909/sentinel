package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"sentinel/internal/app"
	"sentinel/internal/config"
	"sentinel/internal/domain"
	"sentinel/internal/loadbalancer"
	"sentinel/internal/proxy"
	"sentinel/internal/server"
)

func TestServer_ServeHTTPSuccess(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("routed to upstream"))
	}))
	defer upstream.Close()

	u, _ := url.Parse(upstream.URL)
	backend := domain.NewBackend(u, domain.BackendStateHealthy)
	lb, _ := loadbalancer.New(config.RoundRobin)
	svc := &domain.Service{Name: "test-service", Balancer: lb, Backends: []*domain.Backend{backend}}
	route := &domain.Route{Path: "/api", Service: svc}

	rt := &app.Runtime{Routes: []*domain.Route{route}}
	srv := server.New(rt, proxy.New(nil), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "routed to upstream" {
		t.Errorf("expected body 'routed to upstream', got %q", string(body))
	}
}

func TestServer_ServeHTTPNotFound(t *testing.T) {
	srv := server.New(&app.Runtime{}, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found, got %d", res.StatusCode)
	}
}

func TestServer_ServeHTTPNoBackendAvailable(t *testing.T) {
	lb, _ := loadbalancer.New(config.RoundRobin)
	svc := &domain.Service{Name: "empty-service", Balancer: lb, Backends: nil}
	route := &domain.Route{Path: "/empty", Service: svc}

	rt := &app.Runtime{Routes: []*domain.Route{route}}
	srv := server.New(rt, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/empty", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway when no backends available, got %d", res.StatusCode)
	}
}
