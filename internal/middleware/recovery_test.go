package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sentinel/internal/middleware"
)

func TestRecovery_CatchPanic(t *testing.T) {
	handler := middleware.Recovery(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
