package health

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"sentinel/internal/app"
	"sentinel/internal/domain"
)

type backendTracker struct {
	mu                 sync.Mutex
	consecutiveSuccess int
	consecutiveFailure int
}

type Checker struct {
	rt       *app.Runtime
	logger   *slog.Logger
	client   *http.Client
	trackers map[string]*backendTracker
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

func NewChecker(rt *app.Runtime, logger *slog.Logger) *Checker {
	if logger == nil {
		logger = slog.Default()
	}
	return &Checker{
		rt:       rt,
		logger:   logger,
		client:   &http.Client{},
		trackers: make(map[string]*backendTracker),
	}
}

func (c *Checker) getTracker(b *domain.Backend) *backendTracker {
	key := b.URL.String()
	c.mu.Lock()
	defer c.mu.Unlock()
	t, exists := c.trackers[key]
	if !exists {
		t = &backendTracker{}
		c.trackers[key] = t
	}
	return t
}

func (c *Checker) CheckBackend(ctx context.Context, svc *domain.Service, b *domain.Backend) {
	if ctx.Err() != nil {
		return
	}

	targetURL, err := url.JoinPath(b.URL.String(), svc.Health.Path)
	if err != nil {
		c.logger.Error("Invalid health path URL", "service", svc.Name, "error", err)
		c.recordResult(svc, b, false)
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx, svc.Health.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, targetURL, nil)
	if err != nil {
		c.recordResult(svc, b, false)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		c.recordResult(svc, b, false)
		return
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	c.recordResult(svc, b, success)
}

func (c *Checker) recordResult(svc *domain.Service, b *domain.Backend, success bool) {
	t := c.getTracker(b)
	t.mu.Lock()
	defer t.mu.Unlock()

	if success {
		t.consecutiveSuccess++
		t.consecutiveFailure = 0
		if b.GetState() != domain.BackendStateHealthy && t.consecutiveSuccess >= svc.Health.HealthyThreshold {
			c.transitionHealthy(svc, b)
		}
	} else {
		t.consecutiveFailure++
		t.consecutiveSuccess = 0
		if b.GetState() != domain.BackendStateUnhealthy && t.consecutiveFailure >= svc.Health.UnhealthyThreshold {
			c.transitionUnhealthy(svc, b)
		}
	}
}

func (c *Checker) transitionHealthy(svc *domain.Service, b *domain.Backend) {
	oldState := b.GetState()
	b.SetState(domain.BackendStateHealthy)
	c.logger.Info("Backend health state changed", "service", svc.Name, "backend", b.URL.String(), "old_state", oldState.String(), "new_state", "healthy")
}

func (c *Checker) transitionUnhealthy(svc *domain.Service, b *domain.Backend) {
	oldState := b.GetState()
	b.SetState(domain.BackendStateUnhealthy)
	c.logger.Warn("Backend health state changed", "service", svc.Name, "backend", b.URL.String(), "old_state", oldState.String(), "new_state", "unhealthy")
}

func (c *Checker) Start(ctx context.Context) {
	if c.rt == nil || c.rt.Services == nil {
		return
	}
	for _, svc := range c.rt.Services {
		if svc == nil || len(svc.Backends) == 0 {
			continue
		}
		for _, b := range svc.Backends {
			if b == nil {
				continue
			}
			c.wg.Add(1)
			go c.runProbeLoop(ctx, svc, b)
		}
	}
}

func (c *Checker) runProbeLoop(ctx context.Context, svc *domain.Service, b *domain.Backend) {
	defer c.wg.Done()
	if ctx.Err() != nil {
		return
	}
	c.CheckBackend(ctx, svc, b)

	interval := svc.Health.Interval
	if interval <= 0 {
		interval = 10 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.CheckBackend(ctx, svc, b)
		}
	}
}

func (c *Checker) Stop() {
	c.wg.Wait()
}
