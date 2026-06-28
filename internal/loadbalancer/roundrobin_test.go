package loadbalancer_test

import (
	"net/url"
	"testing"

	"sentinel/internal/domain"
	"sentinel/internal/loadbalancer"
)

func TestRoundRobinBalancer_NextBackend(t *testing.T) {
	b := &loadbalancer.RoundRobinBalancer{}

	u1, _ := url.Parse("http://backend1:8080")
	u2, _ := url.Parse("http://backend2:8080")
	u3, _ := url.Parse("http://backend3:8080")

	backends := []*domain.Backend{
		domain.NewBackend(u1, domain.BackendStateHealthy),
		domain.NewBackend(u2, domain.BackendStateHealthy),
		domain.NewBackend(u3, domain.BackendStateHealthy),
	}

	// Verify the first request goes to backend #0
	got, err := b.NextBackend(backends)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.URL.String() != "http://backend1:8080" {
		t.Errorf("expected backend1, got %s", got.URL.String())
	}

	// Second request -> backend #1
	got, err = b.NextBackend(backends)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.URL.String() != "http://backend2:8080" {
		t.Errorf("expected backend2, got %s", got.URL.String())
	}

	// Third request -> backend #2
	got, err = b.NextBackend(backends)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.URL.String() != "http://backend3:8080" {
		t.Errorf("expected backend3, got %s", got.URL.String())
	}

	// Fourth request wraps around -> backend #0
	got, err = b.NextBackend(backends)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.URL.String() != "http://backend1:8080" {
		t.Errorf("expected backend1, got %s", got.URL.String())
	}
}

func TestRoundRobinBalancer_NoBackends(t *testing.T) {
	b := &loadbalancer.RoundRobinBalancer{}

	_, err := b.NextBackend(nil)
	if err == nil {
		t.Errorf("expected error when no backends available, got nil")
	}

	_, err = b.NextBackend([]*domain.Backend{})
	if err == nil {
		t.Errorf("expected error for empty backends slice, got nil")
	}
}
