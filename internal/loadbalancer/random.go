package loadbalancer

import (
	"errors"
	"math/rand/v2"

	"sentinel/internal/domain"
)

type RandomBalancer struct{}

func (r *RandomBalancer) NextBackend(backends []*domain.Backend) (*domain.Backend, error) {
	if len(backends) == 0 {
		return nil, errors.New("no backends available")
	}

	idx := rand.IntN(len(backends))
	return backends[idx], nil
}
