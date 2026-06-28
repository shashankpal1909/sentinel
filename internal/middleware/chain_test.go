package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sentinel/internal/middleware"
)

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

	expected := []string{"mw1-start", "mw2-start", "handler", "mw2-end", "mw1-end"}
	if len(executed) != len(expected) {
		t.Fatalf("expected execution chain %v, got %v", expected, executed)
	}
	for i := range expected {
		if executed[i] != expected[i] {
			t.Errorf("step %d: expected %q, got %q", i, expected[i], executed[i])
		}
	}
}
