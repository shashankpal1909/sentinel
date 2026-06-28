package health_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"sentinel/internal/app"
	"sentinel/internal/domain"
	"sentinel/internal/health"
)

func TestChecker_Thresholds(t *testing.T) {
	var statusCode atomic.Int32
	statusCode.Store(http.StatusInternalServerError)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(int(statusCode.Load()))
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	b := domain.NewBackend(u, domain.BackendStateHealthy)

	svc := &domain.Service{
		Name:               "test-svc",
		Backends:           []*domain.Backend{b},
		HealthPath:         "/healthz",
		HealthTimeout:      1 * time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	}

	checker := health.NewChecker(nil, nil)
	ctx := context.Background()

	// 1st failure: below threshold, should remain healthy
	checker.CheckBackend(ctx, svc, b, nil)
	if b.GetState() != domain.BackendStateHealthy {
		t.Errorf("expected healthy after 1 failure (threshold 2), got %s", b.GetState())
	}

	// 2nd failure: reaches threshold, transitions to unhealthy
	checker.CheckBackend(ctx, svc, b, nil)
	if b.GetState() != domain.BackendStateUnhealthy {
		t.Errorf("expected unhealthy after 2 failures, got %s", b.GetState())
	}

	// Switch server to return 200 OK
	statusCode.Store(http.StatusOK)

	// 1st success: below healthy threshold, should remain unhealthy
	checker.CheckBackend(ctx, svc, b, nil)
	if b.GetState() != domain.BackendStateUnhealthy {
		t.Errorf("expected unhealthy after 1 success (threshold 2), got %s", b.GetState())
	}

	// 2nd success: reaches threshold, restores to healthy
	checker.CheckBackend(ctx, svc, b, nil)
	if b.GetState() != domain.BackendStateHealthy {
		t.Errorf("expected healthy after 2 successes, got %s", b.GetState())
	}
}

func TestChecker_StartStop(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	b := domain.NewBackend(u, domain.BackendStateUnknown)

	svc := &domain.Service{
		Name:               "bg-svc",
		Backends:           []*domain.Backend{b},
		HealthPath:         "/",
		HealthInterval:     10 * time.Millisecond,
		HealthTimeout:      1 * time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	}

	rt := &app.Runtime{
		Services: map[string]*domain.Service{
			"bg-svc": svc,
		},
	}

	checker := health.NewChecker(rt, nil)
	ctx, cancel := context.WithCancel(context.Background())

	checker.Start(ctx)

	// Wait briefly for probe goroutine to run
	time.Sleep(50 * time.Millisecond)
	cancel()
	checker.Stop()

	if b.GetState() != domain.BackendStateHealthy {
		t.Errorf("expected background checker to set healthy state, got %s", b.GetState())
	}
}
