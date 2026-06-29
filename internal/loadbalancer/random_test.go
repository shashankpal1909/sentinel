package loadbalancer_test

import (
	"net/url"
	"testing"

	"sentinel/internal/domain"
	"sentinel/internal/loadbalancer"
)

func TestRandomBalancer_NextBackend(t *testing.T) {
	b := &loadbalancer.RandomBalancer{}

	u1, _ := url.Parse("http://backend1:8080")
	u2, _ := url.Parse("http://backend2:8080")

	backends := []*domain.Backend{
		domain.NewBackend(u1, domain.BackendStateHealthy),
		domain.NewBackend(u2, domain.BackendStateHealthy),
	}

	got, err := b.NextBackend(backends)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got == nil {
		t.Fatalf("expected backend, got nil")
	}
}

func TestRandomBalancer_NoBackends(t *testing.T) {
	b := &loadbalancer.RandomBalancer{}

	_, err := b.NextBackend(nil)
	if err == nil {
		t.Errorf("expected error when no backends available, got nil")
	}
}
