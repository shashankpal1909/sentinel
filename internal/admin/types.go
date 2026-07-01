package admin

import (
	"context"

	"sentinel/internal/app"
	"sentinel/internal/config"
)

// HealthUpdater defines an interface for updating background health checking when runtime changes.
type HealthUpdater interface {
	UpdateRuntime(ctx context.Context, newRt *app.Runtime)
}

type backendResponse struct {
	URL     string `json:"url"`
	State   string `json:"state"`
	Healthy bool   `json:"healthy"`
}

type clusterResponse struct {
	Name        string                    `json:"name"`
	Strategy    config.BalancerStrategy   `json:"strategy"`
	HealthCheck *config.HealthCheckConfig `json:"health_check,omitempty"`
	Backends    []backendResponse         `json:"backends"`
}

type listenerResponse struct {
	Path    string `json:"path"`
	Service string `json:"service"`
}

type runtimeResponse struct {
	Version           uint64 `json:"version"`
	LoadedAt          string `json:"loaded_at"`
	Services          int    `json:"services"`
	Routes            int    `json:"routes"`
	HealthyBackends   int    `json:"healthy_backends"`
	UnhealthyBackends int    `json:"unhealthy_backends"`
}
