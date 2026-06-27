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

	routes := make([]*domain.Route, 0)

	for _, r := range cfg.Routes {
		svcCfg, ok := cfg.Services[r.Service]
		if !ok {
			return nil, fmt.Errorf("service %q not found", r.Service)
		}

		backends := make([]*domain.Backend, 0)
		for _, b := range svcCfg.Backends {
			url, err := url.Parse(b)
			if err != nil {
				return nil, fmt.Errorf("invalid backend URL %q: %v", b, err)
			}
			backends = append(backends, &domain.Backend{
				URL:   url,
				State: domain.BackendStateHealthy,
			})
		}

		routes = append(routes, &domain.Route{
			Path: r.Path,
			Service: &domain.Service{
				Name:     r.Service,
				Strategy: svcCfg.Strategy,
				Balancer: &loadbalancer.RoundRobinBalancer{},
				Backends: backends,
			},
		})
	}

	return &Runtime{
		Routes: routes,
	}, nil
}
