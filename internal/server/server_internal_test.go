package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleRouteUninitialized(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	s.handleRoute(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 Internal Server Error when fields are nil, got %d", res.StatusCode)
	}
}
