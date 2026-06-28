package domain

import (
	"errors"
	"time"
)

type Service struct {
	Name string

	Backends []*Backend
	Balancer Balancer

	HealthPath         string
	HealthInterval     time.Duration
	HealthTimeout      time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
}

func (s *Service) GetHealthyBackends() []*Backend {
	if s == nil {
		return nil
	}
	healthy := make([]*Backend, 0, len(s.Backends))
	for _, b := range s.Backends {
		if b != nil && b.GetState() == BackendStateHealthy {
			healthy = append(healthy, b)
		}
	}
	return healthy
}

func (s *Service) NextBackend() (*Backend, error) {
	if s == nil || s.Balancer == nil {
		return nil, errors.New("service or balancer is nil")
	}
	healthy := s.GetHealthyBackends()
	if len(healthy) == 0 {
		return nil, errors.New("no healthy backends available")
	}
	return s.Balancer.NextBackend(healthy)
}

func (s *Service) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.Name
}
