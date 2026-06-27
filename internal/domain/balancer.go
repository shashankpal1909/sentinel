package domain

type Balancer interface {
	NextBackend(backends []*Backend) (*Backend, error)
}
