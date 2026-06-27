package router_test

import (
	"fmt"
	"testing"

	"sentinel/internal/domain"
	"sentinel/internal/router"
)

func TestRouter_Match(t *testing.T) {
	userSvc := &domain.Service{Name: "user-service"}
	userProfileSvc := &domain.Service{Name: "user-profile-service"}
	rootSvc := &domain.Service{Name: "root-service"}

	routes := []*domain.Route{
		nil, // Ensure nil routes are safely filtered
		{Path: "/", Service: rootSvc},
		{Path: "/users", Service: userSvc},
		{Path: "/users/profile", Service: userProfileSvc},
	}

	r := router.New(routes)

	tests := []struct {
		name      string
		path      string
		wantSvc   string
		wantMatch bool
	}{
		{"exact match longest prefix", "/users/profile", "user-profile-service", true},
		{"subpath match longest prefix", "/users/profile/settings", "user-profile-service", true},
		{"exact match shorter prefix", "/users", "user-service", true},
		{"subpath match shorter prefix", "/users/123", "user-service", true},
		{"prevent false prefix matching", "/usersadmin", "root-service", true},
		{"root match", "/about", "root-service", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, matched := r.Match(tt.path)
			if matched != tt.wantMatch {
				t.Fatalf("Match(%q) matched = %v, want %v", tt.path, matched, tt.wantMatch)
			}
			if matched && svc.Name != tt.wantSvc {
				t.Errorf("Match(%q) service = %q, want %q", tt.path, svc.Name, tt.wantSvc)
			}
		})
	}
}

func TestRouter_NoMatch(t *testing.T) {
	routes := []*domain.Route{
		{Path: "/api", Service: &domain.Service{Name: "api-service"}},
	}
	r := router.New(routes)

	_, matched := r.Match("/web")
	if matched {
		t.Errorf("expected no match for /web when only /api is registered")
	}
}

func TestRouter_NilAndEmptyRoutes(t *testing.T) {
	r := router.New(nil)
	_, matched := r.Match("/users")
	if matched {
		t.Errorf("expected no match on nil router routes")
	}

	r = router.New([]*domain.Route{nil, nil})
	_, matched = r.Match("/users")
	if matched {
		t.Errorf("expected no match when all routes are nil")
	}
}

func BenchmarkRouterMatch(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000}
	svc := &domain.Service{Name: "benchmark-service"}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Routes_%d", size), func(b *testing.B) {
			routes := make([]*domain.Route, 0, size+1)
			for i := 0; i < size; i++ {
				routes = append(routes, &domain.Route{
					Path:    fmt.Sprintf("/api/v1/resource/endpoint_%d", i),
					Service: svc,
				})
			}
			targetPath := "/api/v1/resource/match_target"
			routes = append(routes, &domain.Route{
				Path:    targetPath,
				Service: svc,
			})

			r := router.New(routes)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r.Match(targetPath)
			}
		})
	}
}
