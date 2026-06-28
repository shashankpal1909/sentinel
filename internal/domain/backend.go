package domain

import (
	"net/url"
	"sync/atomic"
)

type BackendState int

const (
	BackendStateUnknown BackendState = iota
	BackendStateHealthy
	BackendStateUnhealthy
)

type Backend struct {
	URL   *url.URL
	state atomic.Int32
}

func NewBackend(u *url.URL, initial BackendState) *Backend {
	b := &Backend{URL: u}
	b.SetState(initial)
	return b
}

func (b *Backend) GetState() BackendState {
	if b == nil {
		return BackendStateUnknown
	}
	return BackendState(b.state.Load())
}

func (b *Backend) SetState(s BackendState) {
	if b != nil {
		b.state.Store(int32(s))
	}
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
	return b.URL.String() + " [" + b.GetState().String() + "]"
}
