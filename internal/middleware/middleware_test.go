package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sentinel/internal/middleware"
)

func TestRequestID_GenerateNew(t *testing.T) {
	var ctxID string
	handler := middleware.RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID = middleware.GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	headerID := res.Header.Get(middleware.RequestIDHeader)
	if headerID == "" {
		t.Errorf("expected X-Request-ID response header to be set, got empty")
	}
	if ctxID != headerID {
		t.Errorf("expected context ID %q to match header ID %q", ctxID, headerID)
	}
}

func TestRequestID_ReuseExisting(t *testing.T) {
	existingID := "custom-trace-id-999"
	var ctxID string
	handler := middleware.RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID = middleware.GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(middleware.RequestIDHeader, existingID)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	headerID := res.Header.Get(middleware.RequestIDHeader)
	if headerID != existingID {
		t.Errorf("expected X-Request-ID response header %q, got %q", existingID, headerID)
	}
	if ctxID != existingID {
		t.Errorf("expected context ID %q, got %q", existingID, ctxID)
	}
}

func TestGetRequestID_EdgeCases(t *testing.T) {
	if id := middleware.GetRequestID(nil); id != "" {
		t.Errorf("expected empty string for nil context, got %q", id)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if id := middleware.GetRequestID(req.Context()); id != "" {
		t.Errorf("expected empty string for context without ID, got %q", id)
	}
}

func TestRecovery_CatchPanic(t *testing.T) {
	handler := middleware.Recovery()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated fatal panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 Internal Server Error after panic recovery, got %d", res.StatusCode)
	}
}

func TestChain(t *testing.T) {
	var executed []string
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executed = append(executed, "mw1-start")
			next.ServeHTTP(w, r)
			executed = append(executed, "mw1-end")
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executed = append(executed, "mw2-start")
			next.ServeHTTP(w, r)
			executed = append(executed, "mw2-end")
		})
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executed = append(executed, "handler")
	})

	chained := middleware.Chain(finalHandler, mw1, mw2)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	chained.ServeHTTP(rec, req)

	expected := []string{"mw2-start", "mw1-start", "handler", "mw1-end", "mw2-end"}
	if len(executed) != len(expected) {
		t.Fatalf("expected execution chain %v, got %v", expected, executed)
	}
	for i := range expected {
		if executed[i] != expected[i] {
			t.Errorf("step %d: expected %q, got %q", i, expected[i], executed[i])
		}
	}
}

func TestLogger_StatusRecord(t *testing.T) {
	handler := middleware.Logger()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected recorder status 404, got %d", rec.Code)
	}
}

func TestLogger_WriteDefaultStatusAndFlush(t *testing.T) {
	handler := middleware.Logger()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected recorder status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", rec.Body.String())
	}
}
