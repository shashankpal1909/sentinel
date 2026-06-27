package app

import (
	"fmt"
	"net/url"

	"sentinel/internal/config"
	"sentinel/internal/domain"
	"sentinel/internal/loadbalancer"
)

func Build(cfg *config.Config) (*Runtime, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

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

		services[name] = &domain.Service{
			Name:     name,
			Strategy: svcCfg.Strategy,
			Balancer: &loadbalancer.RoundRobinBalancer{},
			Backends: backends,
		}
	}

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
	}

	return &Runtime{
		Routes:   routes,
		Services: services,
	}, nil
}
