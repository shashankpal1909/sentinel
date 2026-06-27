package app

import (
	"fmt"
	"strings"

	"sentinel/internal/domain"
)

type Runtime struct {
	Routes []*domain.Route
}

func (r *Runtime) String() string {
	if r == nil || len(r.Routes) == 0 {
		return "Runtime: <empty>"
	}

	var sb strings.Builder
	sb.WriteString("Runtime Configuration:\n")
	for i, route := range r.Routes {
		if route == nil {
			continue
		}
		fmt.Fprintf(&sb, "  [%d] Route: %s\n", i+1, route.Path)
		if route.Service != nil {
			fmt.Fprintf(&sb, "      Service: %s\n", route.Service.Name)
			if route.Service.Strategy != "" {
				fmt.Fprintf(&sb, "      Strategy: %s\n", route.Service.Strategy)
			}
			if len(route.Service.Backends) > 0 {
				sb.WriteString("      Backends:\n")
				for _, b := range route.Service.Backends {
					if b != nil {
						fmt.Fprintf(&sb, "        - %s\n", b.String())
					}
				}
			} else {
				sb.WriteString("      Backends: <none>\n")
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
