package domain

type Service struct {
	Name     string
	Strategy string

	Backends []*Backend
	Balancer Balancer
}

func (s *Service) NextBackend() (*Backend, error) {
	return s.Balancer.NextBackend(s.Backends)
}

func (s *Service) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.Name
}
