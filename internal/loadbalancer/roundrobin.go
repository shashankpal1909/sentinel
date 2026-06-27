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
	idx := r.current.Add(1)

	if len(backends) == 0 {
		return nil, errors.New("no backends available")
	}

	return backends[int(idx%uint32(len(backends)))], nil
}
