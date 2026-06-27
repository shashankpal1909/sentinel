package app

import (
	"fmt"
	"log/slog"
	"net/url"

	"sentinel/internal/config"
	"sentinel/internal/domain"
	"sentinel/internal/loadbalancer"
)

func Build(cfg *config.Config) (*Runtime, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	slog.Debug("Building runtime environment from configuration")

	// Instantiate services and assign load balancing strategies
	services := make(map[string]*domain.Service)
	for name, svcCfg := range cfg.Services {
		backends := make([]*domain.Backend, 0, len(svcCfg.Backends))
		for _, b := range svcCfg.Backends {
			u, err := url.Parse(b)
			if err != nil {
				return nil, fmt.Errorf("invalid backend URL %q for service %q: %w", b, name, err)
			}
			backends = append(backends, &domain.Backend{
				URL:   u,
				State: domain.BackendStateHealthy,
			})
		}

		balancer, err := loadbalancer.New(svcCfg.Strategy)
		if err != nil {
			return nil, fmt.Errorf("failed to create balancer for service %q: %w", name, err)
		}

		services[name] = &domain.Service{
			Name:     name,
			Balancer: balancer,
			Backends: backends,
		}
		slog.Debug("Initialized domain service", "service", name, "backends", len(backends), "strategy", svcCfg.Strategy)
	}

	// Build route mappings linking prefix paths to existing services
	routes := make([]*domain.Route, 0, len(cfg.Routes))
	for _, r := range cfg.Routes {
		svc, ok := services[r.Service]
		if !ok {
			return nil, fmt.Errorf("service %q not found", r.Service)
		}

		routes = append(routes, &domain.Route{
			Path:    r.Path,
			Service: svc,
		})
		slog.Debug("Registered route mapping", "path", r.Path, "service", r.Service)
	}

	slog.Info("Runtime build completed successfully", "services_count", len(services), "routes_count", len(routes))

	return &Runtime{
		Routes:   routes,
		Services: services,
	}, nil
}
