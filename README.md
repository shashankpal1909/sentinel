# Sentinel

Sentinel is a lightweight, high-performance API Gateway built in Go. It provides clean routing, configurable load balancing strategies, structured logging, and robust reverse proxying to manage upstream microservices.

## Features

- **Dynamic Routing**: Prefix-based HTTP routing to upstream services.
- **Load Balancing**: Built-in support for `round-robin` and `random` load balancing strategies.
- **Middleware Pipeline**:
  - **Recovery**: Graceful panic recovery with structured logging and `500 Internal Server Error` responses.
  - **Request Tracing**: Automatic generation and propagation of unique `X-Request-ID` headers across request contexts.
  - **Structured Logging**: Request tracking with latency, status codes, method, path, and remote address via Go's `slog`.
- **Connection Caching**: Efficient reverse proxy pooling to minimize allocation overhead.

## Architecture Pipeline

When an HTTP request arrives at Sentinel, it flows through a structured middleware chain before being routed and reverse-proxied to an upstream backend service:

```text
Client Request
      │
      ▼
┌─────────────────────────────────────────────────────────┐
│ Server Handler Pipeline                                 │
│                                                         │
│  1. Recovery Middleware   (catches panics, 500 Error)   │
│           │                                             │
│           ▼                                             │
│  2. RequestID Middleware  (injects/reuses X-Request-ID) │
│           │                                             │
│           ▼                                             │
│  3. Logger Middleware     (records status & duration)   │
│           │                                             │
│           ▼                                             │
│  4. Router Match          (matches route to service)    │
│           │                                             │
│           ▼                                             │
│  5. Load Balancer         (Round-Robin backend select)  │
│           │                                             │
│           ▼                                             │
│  6. Reverse Proxy         (forwards request & response) │
└───────────┬─────────────────────────────────────────────┘
            │
            ▼
     Upstream Backend
```

For detailed component descriptions, see [ARCHITECTURE.md](file:///Users/shashank/Documents/Projects/sentinel/ARCHITECTURE.md).

## Getting Started

### Prerequisites

- **Go**: 1.26 or newer
- **Docker & Docker Compose**: For containerized orchestration

### Running Locally

1. **Clone the repository**:
   ```bash
   git clone https://github.com/shashankpal1909/sentinel.git
   cd sentinel
   ```

2. **Run with default configuration**:
   ```bash
   CONFIG_PATH=example.gateway.yaml make run
   ```

### Running with Docker Compose

To launch Sentinel alongside mock upstream services (echo backends):

```bash
make docker-run
# Or directly via docker compose:
docker compose up --build
```

## Configuration

Sentinel uses YAML configuration files to define services, backends, and routes. See [example.gateway.yaml](file:///Users/shashank/Documents/Projects/sentinel/example.gateway.yaml) for a complete reference:

```yaml
server:
  port: 8080

services:
  auth-service:
    strategy: round-robin
    backends:
      - http://localhost:8001
      - http://localhost:8002

routes:
  - path: /login
    service: auth-service
```

## Development & CI

Sentinel includes a full suite of unit and integration tests with **100% statement coverage** across core modules.

Run the CI verification suite locally using `make`:

```bash
# Run all CI checks (fmt, vet, test, race detector, code coverage)
make ci
```
