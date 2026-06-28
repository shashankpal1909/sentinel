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
	trackers map[*domain.Backend]*backendTracker
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
		trackers: make(map[*domain.Backend]*backendTracker),
	}
}

func (c *Checker) getTracker(b *domain.Backend) *backendTracker {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, exists := c.trackers[b]
	if !exists {
		t = &backendTracker{}
		c.trackers[b] = t
	}
	return t
}

func (c *Checker) CheckBackend(ctx context.Context, svc *domain.Service, b *domain.Backend, client *http.Client) {
	if ctx.Err() != nil {
		return
	}
	if client == nil {
		client = &http.Client{Timeout: svc.HealthTimeout}
	}

	targetURL, err := url.JoinPath(b.URL.String(), svc.HealthPath)
	if err != nil {
		c.logger.Error("Invalid health path URL", "service", svc.Name, "error", err)
		c.recordResult(svc, b, false)
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx, svc.HealthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, targetURL, nil)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		c.recordResult(svc, b, false)
		return
	}

	resp, err := client.Do(req)
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

	healthyThresh := svc.HealthyThreshold
	if healthyThresh <= 0 {
		healthyThresh = 1
	}
	unhealthyThresh := svc.UnhealthyThreshold
	if unhealthyThresh <= 0 {
		unhealthyThresh = 1
	}

	if success {
		t.consecutiveSuccess++
		t.consecutiveFailure = 0
		if b.GetState() != domain.BackendStateHealthy && t.consecutiveSuccess >= healthyThresh {
			b.SetState(domain.BackendStateHealthy)
			c.logger.Info("Backend restored to healthy state", "service", svc.Name, "backend", b.URL.String())
		}
	} else {
		t.consecutiveFailure++
		t.consecutiveSuccess = 0
		if b.GetState() != domain.BackendStateUnhealthy && t.consecutiveFailure >= unhealthyThresh {
			b.SetState(domain.BackendStateUnhealthy)
			c.logger.Warn("Backend marked unhealthy", "service", svc.Name, "backend", b.URL.String())
		}
	}
}

func (c *Checker) Start(ctx context.Context) {
	if c.rt == nil || c.rt.Services == nil {
		return
	}
	for _, svc := range c.rt.Services {
		if svc == nil || len(svc.Backends) == 0 {
			continue
		}
		client := &http.Client{Timeout: svc.HealthTimeout}
		for _, b := range svc.Backends {
			if b == nil {
				continue
			}
			c.wg.Add(1)
			go func(service *domain.Service, backend *domain.Backend) {
				defer c.wg.Done()
				if ctx.Err() != nil {
					return
				}
				c.CheckBackend(ctx, service, backend, client)

				interval := service.HealthInterval
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
						c.CheckBackend(ctx, service, backend, client)
					}
				}
			}(svc, b)
		}
	}
}

func (c *Checker) Stop() {
	c.wg.Wait()
}
