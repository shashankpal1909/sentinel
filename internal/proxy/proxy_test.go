package proxy_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"sentinel/internal/domain"
	"sentinel/internal/proxy"
)

func TestProxy_ForwardSuccess(t *testing.T) {
	// Start mock upstream server echoing back request path and header
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Header") != "sentinel" {
			t.Errorf("expected header X-Test-Header=sentinel, got %s", r.Header.Get("X-Test-Header"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("echo from upstream: " + r.URL.Path))
	}))
	defer upstream.Close()

	u, _ := url.Parse(upstream.URL)
	backend := &domain.Backend{URL: u, State: domain.BackendStateHealthy}

	p := proxy.New()
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	req.Header.Set("X-Test-Header", "sentinel")
	rec := httptest.NewRecorder()

	p.Forward(backend, rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	expected := "echo from upstream: /test/path"
	if string(body) != expected {
		t.Errorf("expected body %q, got %q", expected, string(body))
	}
}

func TestProxy_ForwardNilBackend(t *testing.T) {
	p := proxy.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	p.Forward(nil, rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway for nil backend, got %d", res.StatusCode)
	}
}

func TestProxy_ForwardUnreachableBackend(t *testing.T) {
	// Point backend to closed local port
	u, _ := url.Parse("http://127.0.0.1:0")
	backend := &domain.Backend{URL: u, State: domain.BackendStateUnhealthy}

	p := proxy.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	p.Forward(backend, rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway for unreachable backend, got %d", res.StatusCode)
	}
}
