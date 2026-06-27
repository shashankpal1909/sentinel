package router

import (
	"sort"
	"strings"

	"sentinel/internal/domain"
)

type Router struct {
	routes []*domain.Route
}

func New(routes []*domain.Route) *Router {
	cleanRoutes := make([]*domain.Route, 0, len(routes))
	for _, r := range routes {
		if r != nil {
			cleanRoutes = append(cleanRoutes, r)
		}
	}

	sort.Slice(cleanRoutes, func(i, j int) bool {
		return len(cleanRoutes[i].Path) > len(cleanRoutes[j].Path)
	})

	return &Router{routes: cleanRoutes}
}

func (r *Router) Match(path string) (*domain.Service, bool) {
	for _, route := range r.routes {
		if matchPath(path, route.Path) {
			return route.Service, true
		}
	}

	return nil, false
}

func matchPath(reqPath, routePath string) bool {
	if routePath == "/" {
		return strings.HasPrefix(reqPath, "/")
	}
	cleanRoute := strings.TrimRight(routePath, "/")
	if reqPath == cleanRoute {
		return true
	}
	return strings.HasPrefix(reqPath, cleanRoute+"/")
}
