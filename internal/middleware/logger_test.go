package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sentinel/internal/middleware"
)

func TestLogger_StatusRecord(t *testing.T) {
	handler := middleware.Logger(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	handler := middleware.Logger(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
