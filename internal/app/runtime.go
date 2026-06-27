package app

import (
	"fmt"
	"io"
	"strings"

	"sentinel/internal/domain"
)

type Runtime struct {
	Routes   []*domain.Route
	Services map[string]*domain.Service
}

func (r *Runtime) Dump(w io.Writer) {
	if r == nil || (len(r.Routes) == 0 && len(r.Services) == 0) {
		fmt.Fprint(w, "Runtime: <empty>")
		return
	}

	fmt.Fprintln(w, "Runtime Configuration:")
	if len(r.Services) > 0 {
		fmt.Fprintln(w, "  Services:")
		for _, svc := range r.Services {
			if svc == nil {
				continue
			}
			strategy := svc.Strategy
			if strategy == "" {
				strategy = "default"
			}
			fmt.Fprintf(w, "    - %s (Strategy: %s)\n", svc.Name, strategy)
			if len(svc.Backends) > 0 {
				for _, b := range svc.Backends {
					if b != nil {
						fmt.Fprintf(w, "        Backend: %s\n", b.String())
					}
				}
			} else {
				fmt.Fprintln(w, "        Backend: <none>")
			}
		}
	}

	if len(r.Routes) > 0 {
		fmt.Fprintln(w, "  Routes:")
		for i, route := range r.Routes {
			if route == nil {
				continue
			}
			svcName := "<nil>"
			if route.Service != nil {
				svcName = route.Service.Name
			}
			fmt.Fprintf(w, "    [%d] %s -> %s\n", i+1, route.Path, svcName)
		}
	}
}

func (r *Runtime) String() string {
	var sb strings.Builder
	r.Dump(&sb)
	return strings.TrimRight(sb.String(), "\n")
}
