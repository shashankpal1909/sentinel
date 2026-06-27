package domain

import "net/url"

type BackendState int

const (
	BackendStateUnknown BackendState = iota
	BackendStateHealthy
	BackendStateUnhealthy
)

type Backend struct {
	URL   *url.URL
	State BackendState
}

func (s BackendState) String() string {
	switch s {
	case BackendStateHealthy:
		return "healthy"
	case BackendStateUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

func (b *Backend) String() string {
	if b == nil || b.URL == nil {
		return "<nil>"
	}
	return b.URL.String() + " [" + b.State.String() + "]"
}
