package domain_test

import (
	"net/url"
	"sync"
	"testing"

	"sentinel/internal/domain"
)

type mockBalancer struct {
	backend  *domain.Backend
	received []*domain.Backend
}

func (m *mockBalancer) NextBackend(backends []*domain.Backend) (*domain.Backend, error) {
	m.received = backends
	return m.backend, nil
}

func TestService_NextBackendAndString(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080")
	b := domain.NewBackend(u, domain.BackendStateHealthy)

	svc := &domain.Service{
		Name:     "test-svc",
		Backends: []*domain.Backend{b},
		Balancer: &mockBalancer{backend: b},
	}

	got, err := svc.NextBackend()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b {
		t.Errorf("expected backend %v, got %v", b, got)
	}

	if svc.String() != "test-svc" {
		t.Errorf("expected test-svc, got %s", svc.String())
	}

	var nilSvc *domain.Service
	if nilSvc.String() != "<nil>" {
		t.Errorf("expected <nil> for nil service string, got %s", nilSvc.String())
	}
}

func TestService_HealthFiltering(t *testing.T) {
	u1, _ := url.Parse("http://localhost:8081")
	u2, _ := url.Parse("http://localhost:8082")
	b1 := domain.NewBackend(u1, domain.BackendStateUnhealthy)
	b2 := domain.NewBackend(u2, domain.BackendStateHealthy)

	mb := &mockBalancer{backend: b2}
	svc := &domain.Service{
		Name:     "test-svc",
		Backends: []*domain.Backend{b1, b2},
		Balancer: mb,
	}

	healthy := svc.GetHealthyBackends()
	if len(healthy) != 1 || healthy[0] != b2 {
		t.Fatalf("expected only b2 healthy, got %v", healthy)
	}

	got, err := svc.NextBackend()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != b2 {
		t.Errorf("expected b2, got %v", got)
	}
	if len(mb.received) != 1 || mb.received[0] != b2 {
		t.Errorf("balancer received unfiltered or wrong backends: %v", mb.received)
	}

	// Test all unhealthy
	b2.SetState(domain.BackendStateUnhealthy)
	_, err = svc.NextBackend()
	if err == nil {
		t.Errorf("expected error when all backends are unhealthy, got nil")
	}

	var nilSvc *domain.Service
	if nilSvc.GetHealthyBackends() != nil {
		t.Errorf("expected nil for nil service GetHealthyBackends")
	}
}

func TestRoute_String(t *testing.T) {
	var nilRoute *domain.Route
	if nilRoute.String() != "<nil>" {
		t.Errorf("expected <nil>, got %s", nilRoute.String())
	}

	r := &domain.Route{Path: "/test", Service: nil}
	if r.String() != "/test -> <nil>" {
		t.Errorf("expected '/test -> <nil>', got %s", r.String())
	}

	r.Service = &domain.Service{Name: "my-svc"}
	if r.String() != "/test -> my-svc" {
		t.Errorf("expected '/test -> my-svc', got %s", r.String())
	}
}

func TestBackendAndState_String(t *testing.T) {
	var nilBackend *domain.Backend
	if nilBackend.String() != "<nil>" {
		t.Errorf("expected <nil>, got %s", nilBackend.String())
	}

	u, _ := url.Parse("http://localhost:9000")
	b := domain.NewBackend(u, domain.BackendStateHealthy)
	if b.String() != "http://localhost:9000 [healthy]" {
		t.Errorf("expected healthy representation, got %s", b.String())
	}

	b.SetState(domain.BackendStateUnhealthy)
	if b.String() != "http://localhost:9000 [unhealthy]" {
		t.Errorf("expected unhealthy representation, got %s", b.String())
	}

	b.SetState(domain.BackendStateUnknown)
	if b.String() != "http://localhost:9000 [unknown]" {
		t.Errorf("expected unknown representation, got %s", b.String())
	}

	b.SetState(domain.BackendState(999))
	if b.GetState().String() != "unknown" {
		t.Errorf("expected unknown state for invalid int, got %s", b.GetState().String())
	}
}

func TestBackend_ConcurrentState(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080")
	b := domain.NewBackend(u, domain.BackendStateHealthy)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			b.SetState(domain.BackendStateUnhealthy)
			b.SetState(domain.BackendStateHealthy)
		}()
		go func() {
			defer wg.Done()
			_ = b.GetState()
			_ = b.String()
		}()
	}
	wg.Wait()
}
