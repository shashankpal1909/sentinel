package app

import (
	"fmt"
	"time"

	"sentinel/internal/config"
	"sentinel/internal/router"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Build(cfg *config.Config, version uint64) (*Snapshot, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if err := config.Validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	rt, err := Build(cfg)
	if err != nil {
		return nil, fmt.Errorf("runtime build failed: %w", err)
	}
	rtr := router.New(rt.Routes)
	return &Snapshot{
		Runtime:  rt,
		Router:   rtr,
		Version:  version,
		LoadedAt: time.Now(),
	}, nil
}

func (l *Loader) Load(path string, version uint64) (*config.Config, *Snapshot, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config file: %w", err)
	}
	snap, err := l.Build(cfg, version)
	if err != nil {
		return nil, nil, err
	}
	return cfg, snap, nil
}
