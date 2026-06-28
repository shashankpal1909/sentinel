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

func TestProxy_ForwardSuccessAndCaching(t *testing.T) {
	hits := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.Header.Get("X-Test-Header") != "sentinel" {
			t.Errorf("expected header X-Test-Header=sentinel, got %s", r.Header.Get("X-Test-Header"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("echo from upstream: " + r.URL.Path))
	}))
	defer upstream.Close()

	u, _ := url.Parse(upstream.URL)
	backend := domain.NewBackend(u, domain.BackendStateHealthy)

	p := proxy.New(nil)

	// Send multiple sequential requests to verify reverse proxy reuse/caching
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
		req.Header.Set("X-Test-Header", "sentinel")
		rec := httptest.NewRecorder()

		p.Forward(rec, req, backend)

		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 OK on request %d, got %d", i, res.StatusCode)
		}

		body, _ := io.ReadAll(res.Body)
		expected := "echo from upstream: /test/path"
		if string(body) != expected {
			t.Errorf("expected body %q on request %d, got %q", expected, i, string(body))
		}
	}

	if hits != 3 {
		t.Errorf("expected 3 hits on upstream server, got %d", hits)
	}
}

func TestProxy_ForwardNilBackend(t *testing.T) {
	p := proxy.New(nil)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	p.Forward(rec, req, nil)

	res := rec.Result()
	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway for nil backend, got %d", res.StatusCode)
	}
}

func TestProxy_ForwardUnreachableBackend(t *testing.T) {
	u, _ := url.Parse("http://127.0.0.1:0")
	backend := domain.NewBackend(u, domain.BackendStateUnhealthy)

	p := proxy.New(nil)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	p.Forward(rec, req, backend)

	res := rec.Result()
	if res.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway for unreachable backend, got %d", res.StatusCode)
	}
}
