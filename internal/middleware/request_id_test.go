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
