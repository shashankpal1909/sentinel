package loadbalancer

import (
	"fmt"

	"sentinel/internal/config"
	"sentinel/internal/domain"
)

func New(strategy config.BalancerStrategy) (domain.Balancer, error) {
	switch strategy {
	case config.RoundRobin, "":
		return &RoundRobinBalancer{}, nil
	case config.Random:
		return &RandomBalancer{}, nil
	default:
		return nil, fmt.Errorf("unknown balancer strategy %q", strategy)
	}
}
