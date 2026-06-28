package config

import (
	"fmt"
	"log/slog"
	"os"

	"go.yaml.in/yaml/v4"
)

func Load(paths ...string) (*Config, error) {
	path := "gateway.yaml"
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		path = env
	}
	if len(paths) > 0 && paths[0] != "" {
		path = paths[0]
	}

	slog.Info("Loading configuration", "path", path)

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML config from %q: %w", path, err)
	}

	slog.Debug("Configuration loaded successfully", "services", len(cfg.Services), "routes", len(cfg.Routes))
	return &cfg, nil
}
