package domain

type Route struct {
	Path string

	Service *Service
}

func (r *Route) String() string {
	if r == nil {
		return "<nil>"
	}
	if r.Service == nil {
		return r.Path + " -> <nil>"
	}
	return r.Path + " -> " + r.Service.Name
}
