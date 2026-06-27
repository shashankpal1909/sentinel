package config

import (
	"os"

	"go.yaml.in/yaml/v4"
)

func Load(paths ...string) (*Config, error) {
	path := "gateway.yaml"
	if len(paths) > 0 && paths[0] != "" {
		path = paths[0]
	}

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
