package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"sentinel/internal/config"
)

// HealthUpdater defines an interface for updating background health checking when runtime changes.
type HealthUpdater interface {
	UpdateRuntime(ctx context.Context, newRt *Runtime)
}

type Manager struct {
	configPath string
	logger     *slog.Logger
	cfg        atomic.Pointer[config.Config]
	rt         atomic.Pointer[Runtime]
	health     HealthUpdater
	ctx        context.Context
	mu         sync.Mutex
}

func NewManager(initialCfg *config.Config, configPath string, logger *slog.Logger) (*Manager, error) {
	if initialCfg == nil {
		return nil, errors.New("initial config cannot be nil")
	}
	if err := config.Validate(initialCfg); err != nil {
		return nil, fmt.Errorf("invalid initial config: %w", err)
	}
	if logger == nil {
		logger = slog.Default()
	}
	rt, err := Build(initialCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build runtime: %w", err)
	}

	m := &Manager{
		configPath: configPath,
		logger:     logger,
	}
	m.cfg.Store(initialCfg)
	m.rt.Store(rt)
	return m, nil
}

func (m *Manager) SetHealthUpdater(ctx context.Context, h HealthUpdater) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ctx = ctx
	m.health = h
}

func (m *Manager) GetConfig() *config.Config {
	return m.cfg.Load()
}

func (m *Manager) GetRuntime() *Runtime {
	return m.rt.Load()
}

func (m *Manager) Apply(newCfg *config.Config) error {
	if newCfg == nil {
		return errors.New("new config cannot be nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := config.Validate(newCfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	newRt, err := Build(newCfg)
	if err != nil {
		return fmt.Errorf("runtime build failed: %w", err)
	}

	// Preserve atomic health states for existing backends
	oldRt := m.rt.Load()
	if oldRt != nil && oldRt.Services != nil {
		for svcName, newSvc := range newRt.Services {
			oldSvc, exists := oldRt.Services[svcName]
			if !exists || oldSvc == nil {
				continue
			}
			for _, newB := range newSvc.Backends {
				for _, oldB := range oldSvc.Backends {
					if newB.URL != nil && oldB.URL != nil && newB.URL.String() == oldB.URL.String() {
						newB.SetState(oldB.GetState())
						break
					}
				}
			}
		}
	}

	m.cfg.Store(newCfg)
	m.rt.Store(newRt)

	if m.health != nil && m.ctx != nil {
		m.health.UpdateRuntime(m.ctx, newRt)
	}

	m.logger.Info("Configuration applied successfully via hot reload", "services", len(newRt.Services), "routes", len(newRt.Routes))
	return nil
}

func (m *Manager) ReloadFromDisk() error {
	if m.configPath == "" {
		return errors.New("no configuration file path set for reload")
	}
	m.logger.Info("Reloading configuration from disk", "path", m.configPath)
	newCfg, err := config.Load(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration file: %w", err)
	}
	return m.Apply(newCfg)
}
