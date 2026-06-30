package app

import (
	"errors"
	"sync"
	"sync/atomic"

	"sentinel/internal/config"
	"sentinel/internal/domain"
)

type Manager struct {
	cfg  atomic.Pointer[config.Config]
	snap atomic.Pointer[Snapshot]
	mu   sync.Mutex
}

func NewManager(initialSnap *Snapshot, initialCfg *config.Config) (*Manager, error) {
	if initialSnap == nil {
		return nil, errors.New("initial snapshot cannot be nil")
	}
	if initialCfg == nil {
		return nil, errors.New("initial config cannot be nil")
	}
	m := &Manager{}
	m.snap.Store(initialSnap)
	m.cfg.Store(initialCfg)
	return m, nil
}

func (m *Manager) CurrentConfig() *config.Config {
	return m.cfg.Load()
}

// GetConfig is preserved for backwards compatibility with existing callers.
func (m *Manager) GetConfig() *config.Config {
	return m.CurrentConfig()
}

func (m *Manager) Current() *Snapshot {
	return m.snap.Load()
}

// GetRuntime is preserved for backwards compatibility pointing to Current().Runtime.
func (m *Manager) GetRuntime() *Runtime {
	s := m.Current()
	if s == nil {
		return nil
	}
	return s.Runtime
}

func (m *Manager) Replace(newSnap *Snapshot, newCfg ...*config.Config) {
	if newSnap == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	oldSnap := m.snap.Load()
	if oldSnap != nil && oldSnap.Runtime != nil && newSnap.Runtime != nil {
		states := make(map[string]domain.BackendState)
		for _, svc := range oldSnap.Runtime.Services {
			if svc == nil {
				continue
			}
			for _, b := range svc.Backends {
				if b != nil && b.URL != nil {
					states[b.URL.String()] = b.GetState()
				}
			}
		}
		for _, svc := range newSnap.Runtime.Services {
			if svc == nil {
				continue
			}
			for _, b := range svc.Backends {
				if b != nil && b.URL != nil {
					if st, ok := states[b.URL.String()]; ok {
						b.SetState(st)
					}
				}
			}
		}
	}

	m.snap.Store(newSnap)
	if len(newCfg) > 0 && newCfg[0] != nil {
		m.cfg.Store(newCfg[0])
	}
}
