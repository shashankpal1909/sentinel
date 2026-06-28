# Sentinel Architecture

Sentinel is a lightweight, high-performance API Gateway designed with clean architectural boundaries, modular pipeline execution, and fail-fast invariants.

## Package Overview

```text
sentinel/
├── cmd/sentinel/         # Application entrypoint
├── internal/
│   ├── app/              # Runtime initializer and domain assembler
│   ├── config/           # YAML configuration loader and validator
│   ├── domain/           # Core domain models and state entities
│   ├── loadbalancer/     # Load balancing strategies (Round-Robin, Random)
│   ├── logger/           # Structured logging setup (slog)
│   ├── middleware/       # HTTP middleware pipeline and chain orchestration
│   ├── proxy/            # Reverse proxy wrapper and connection cache
│   ├── router/           # Trie/Prefix-based route matching engine
│   └── server/           # HTTP Gateway server orchestrator
└── integration/          # End-to-end gateway integration tests
```

## Core Architectural Components

### 1. Configuration & Loader (`internal/config`)
Loads and parses YAML gateway definitions (`gateway.yaml` or via `CONFIG_PATH`). Validates that all defined routes point to existing services and all backends have valid HTTP/HTTPS URLs before the application starts.

### 2. Domain Models (`internal/domain`)
Defines the core entities:
- `Backend`: Represents an upstream target URL and its health state (`Healthy`, `Unhealthy`, `Draining`).
- `Service`: Groups multiple backends under a logical name assigned to a load balancing strategy.
- `Route`: Maps an incoming URL path prefix to a specific `Service`.

### 3. Load Balancer (`internal/loadbalancer`)
Provides thread-safe backend selection algorithms:
- **Round-Robin**: Sequentially cycles through healthy backends using atomic counters.
- **Random**: Selects a healthy backend uniformly at random using fast math/rand generators.

### 4. Router Engine (`internal/router`)
Matches incoming request paths against configured routes. Supports prefix matching to route traffic to the appropriate upstream service.

### 5. Middleware Pipeline (`internal/middleware`)
Intercepts HTTP requests before they reach the core routing logic. Structured as a backward-wrapping chain so execution happens strictly left-to-right:
1. **Recovery**: Captures runtime panics, logs stack traces via `slog`, and returns clean `500 Internal Server Error` responses.
2. **RequestID**: Checks for incoming `X-Request-ID` headers or generates UUIDv4 tokens. Propagates the ID to both request/response headers and `context.Context`.
3. **Logger**: Wraps `http.ResponseWriter` to record status codes and response latencies, emitting structured logs upon completion.

### 6. Reverse Proxy (`internal/proxy`)
Wraps Go's `httputil.ReverseProxy`. Caches proxy instances per target URL to eliminate repeated allocations and customizes error handling to return clean `502 Bad Gateway` responses when upstream targets are unreachable.

### 7. Gateway Server (`internal/server`)
Orchestrates the router, proxy, and middleware chain. Enforces fail-fast constructor invariants (`server.New`), eliminating redundant runtime checks during request servicing.
