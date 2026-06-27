package loadbalancer

import (
	"errors"
	"sync/atomic"

	"sentinel/internal/domain"
)

type RoundRobinBalancer struct {
	current atomic.Uint32
}

func (r *RoundRobinBalancer) NextBackend(backends []*domain.Backend) (*domain.Backend, error) {
	if len(backends) == 0 {
		return nil, errors.New("no backends available")
	}

	idx := r.current.Add(1) - 1
	return backends[int(idx%uint32(len(backends)))], nil
}
