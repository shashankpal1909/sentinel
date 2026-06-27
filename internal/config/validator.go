package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func Validate(cfg *Config) error {
	if cfg == nil {
		return errors.New("config cannot be nil")
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if len(cfg.Services) == 0 {
		return errors.New("no services defined in configuration")
	}

	for name, svc := range cfg.Services {
		if strings.TrimSpace(name) == "" {
			return errors.New("service name cannot be empty")
		}
		if len(svc.Backends) == 0 {
			return fmt.Errorf("service %q has empty backend list", name)
		}
		for _, b := range svc.Backends {
			if strings.TrimSpace(b) == "" {
				return fmt.Errorf("service %q has empty backend URL", name)
			}
			u, err := url.ParseRequestURI(b)
			if err != nil {
				return fmt.Errorf("service %q has invalid backend URL %q: %w", name, b, err)
			}
			if u.Scheme == "" || u.Host == "" {
				return fmt.Errorf("service %q has invalid backend URL %q: missing scheme or host", name, b)
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				return fmt.Errorf("service %q backend URL %q must use http or https scheme", name, b)
			}
		}
	}

	seenPaths := make(map[string]bool)
	for _, r := range cfg.Routes {
		if strings.TrimSpace(r.Path) == "" {
			return errors.New("route path cannot be empty")
		}
		if !strings.HasPrefix(r.Path, "/") {
			return fmt.Errorf("route path %q must start with '/'", r.Path)
		}
		if seenPaths[r.Path] {
			return fmt.Errorf("duplicate route path: %q", r.Path)
		}
		seenPaths[r.Path] = true

		if _, ok := cfg.Services[r.Service]; !ok {
			return fmt.Errorf("route %q references missing service %q", r.Path, r.Service)
		}
	}

	return nil
}
