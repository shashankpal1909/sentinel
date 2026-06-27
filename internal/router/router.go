package router

import (
	"log/slog"
	"sort"
	"strings"

	"sentinel/internal/domain"
)

type Router struct {
	routes []*domain.Route
}

func New(routes []*domain.Route) *Router {
	// Filter nil entries to prevent dereference panics during matching
	cleanRoutes := make([]*domain.Route, 0, len(routes))
	for _, r := range routes {
		if r != nil {
			cleanRoutes = append(cleanRoutes, r)
		}
	}

	// Sort longest path prefixes first so specific subpaths take priority over general prefixes
	sort.Slice(cleanRoutes, func(i, j int) bool {
		return len(cleanRoutes[i].Path) > len(cleanRoutes[j].Path)
	})

	slog.Debug("Initialized router table", "active_routes", len(cleanRoutes))
	return &Router{routes: cleanRoutes}
}

func (r *Router) Match(path string) (*domain.Service, bool) {
	for _, route := range r.routes {
		if matchPath(path, route.Path) {
			slog.Debug("Route match found", "request_path", path, "matched_prefix", route.Path, "service", route.Service.Name)
			return route.Service, true
		}
	}

	slog.Debug("No route match found", "request_path", path)
	return nil, false
}

func matchPath(reqPath, routePath string) bool {
	if routePath == "/" {
		return strings.HasPrefix(reqPath, "/")
	}
	// Prevent false matches across word boundaries (e.g. /users matching /usersadmin)
	cleanRoute := strings.TrimRight(routePath, "/")
	if reqPath == cleanRoute {
		return true
	}
	return strings.HasPrefix(reqPath, cleanRoute+"/")
}
