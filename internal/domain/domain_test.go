package domain_test

import (
	"net/url"
	"testing"

	"sentinel/internal/domain"
)

type mockBalancer struct {
	backend *domain.Backend
}

func (m *mockBalancer) NextBackend(backends []*domain.Backend) (*domain.Backend, error) {
	return m.backend, nil
}

func TestService_NextBackendAndString(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080")
	b := &domain.Backend{URL: u, State: domain.BackendStateHealthy}

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
	b := &domain.Backend{URL: u, State: domain.BackendStateHealthy}
	if b.String() != "http://localhost:9000 [healthy]" {
		t.Errorf("expected healthy representation, got %s", b.String())
	}

	b.State = domain.BackendStateUnhealthy
	if b.String() != "http://localhost:9000 [unhealthy]" {
		t.Errorf("expected unhealthy representation, got %s", b.String())
	}

	b.State = domain.BackendStateUnknown
	if b.String() != "http://localhost:9000 [unknown]" {
		t.Errorf("expected unknown representation, got %s", b.String())
	}

	b.State = domain.BackendState(999)
	if b.State.String() != "unknown" {
		t.Errorf("expected unknown state for invalid int, got %s", b.State.String())
	}
}
