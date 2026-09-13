# Load Balancer in Go

A simple HTTP load balancer built with **Go**, using **Python backend servers** for local testing.

This project was built to learn how load balancers work under the hood — reverse proxying, health checks, and the tradeoffs between different load-balancing algorithms.

## Overview

The load balancer receives HTTP requests and forwards them to healthy backend servers. Load-balancing strategies are implemented behind a common `Algorithm` interface, so the distribution logic is fully decoupled from the proxying and health-check machinery. Which strategy runs is chosen via `config.json` — no code changes needed to switch algorithms.

```text
Client
   |
   v
Go Load Balancer (:8081)
   |
   +----> Python Backend (:5001)
   |
   +----> Python Backend (:5002)
   |
   +----> Python Backend (:5003)
   |
   +----> Python Backend (:5004)
   |
   +----> Python Backend (:5005)
```

## Project Structure

```text
cmd/load-balancer/    Entry point (main.go) and config.json
internal/balancer/    Server model, config struct, health checks, proxy pooling
internal/algorithms/  Load-balancing strategies (Algorithm implementations)
server/server.py      Minimal Python backend used for local testing
```

## Features

* Pluggable load-balancing strategies via a shared `Algorithm` interface, selected in `config.json`
* **Round Robin** — cycles through healthy servers in order
* **Weighted Round Robin** — a smooth weighted round-robin: each server accumulates its weight every round, the highest accumulator is chosen, then reduced by the total weight. Weights are randomly assigned (a unique value from 1 to the number of servers) rather than derived from real server metrics like CPU or memory
* **Least Connections** — routes to the healthy server with the fewest active connections, tracked via a mutex-guarded counter on each server
* **IP Hash** — hashes a routing key (the client IP, or `X-Forwarded-For`/`X-Real-IP` if present) with FNV-1a and mods by the number of currently healthy servers, so the same client sticks to the same backend as long as the healthy server count doesn't change
* Periodic backend health checks (HTTP HEAD on an interval), with automatic skipping of unhealthy servers
* Reverse-proxy connection pooling (`ProxyPool`) — one `httputil.ReverseProxy` is built per server and reused, instead of rebuilding it on every request
* `X-Forwarded-Server` response header naming which backend handled each request

## Configuration

Everything is configured in `cmd/load-balancer/config.json`:

```json
{
  "port": ":8081",
  "healthCheckInterval": "2s",
  "algorithm": "weighted-round-robin",
  "servers": [
    "http://localhost:5001",
    "http://localhost:5002",
    "http://localhost:5003",
    "http://localhost:5004",
    "http://localhost:5005"
  ]
}
```

`algorithm` accepts one of: `round-robin`, `weighted-round-robin`, `least-connections`, `ip-hash`. An unrecognized value causes the load balancer to exit at startup.

## Running the Project

### 1. Start the Python backend servers

Open five terminal windows and run:

```bash
PORT=5001 python3 server/server.py
PORT=5002 python3 server/server.py
PORT=5003 python3 server/server.py
PORT=5004 python3 server/server.py
PORT=5005 python3 server/server.py
```

Each server listens on its own port.

### 2. Start the Go load balancer

`main.go` loads `config.json` from its own working directory, so run it from inside `cmd/load-balancer`:

```bash
cd cmd/load-balancer
go run .
```

The load balancer starts on:

```text
http://localhost:8081
```

### 3. Send requests

```bash
curl http://localhost:8081
```

Example response:

```text
Response from Python server 5001
```

Send multiple requests to observe the load-balancing distribution:

```bash
for i in {1..10}; do curl http://localhost:8081; done
```

Each response also includes an `X-Forwarded-Server` header naming the backend that handled the request.

## How It Works

1. The load balancer reads `config.json` and constructs the selected `Algorithm`.
2. Each configured server URL becomes a `Server`, assigned a random unique weight (1 to the number of servers) and marked healthy by default.
3. A goroutine is started per server, periodically sending an HTTP HEAD request and marking the server unhealthy if it fails or doesn't return `200`.
4. For every incoming request, the client's routing key is resolved (`X-Forwarded-For` → `X-Real-IP` → `RemoteAddr`), and the active `Algorithm` selects the next healthy server using that key (Round Robin and Least Connections ignore it; IP Hash uses it).
5. The request is forwarded through a pooled `httputil.ReverseProxy` for that server, with the server's active-connection count incremented before the call and decremented after, via `defer`.
6. If no healthy servers are available, the load balancer returns `503 Service Unavailable`.

## Learning Goals

This project was built to understand:

* HTTP request forwarding and reverse proxies
* Load-balancing strategies — round robin, weighted round robin, least connections, IP hash — and the tradeoffs between them
* Interface-based design for swapping algorithms without touching the proxying/health-check code
* Health checks running as independent goroutines
* Goroutines, concurrency, and protecting shared state with mutexes
* Basic Go networking (`net/http`, `net/http/httputil`)

## Known Limitations

This is a learning project, not a production load balancer. Some deliberate simplifications:

* Weighted Round Robin's weights are randomly assigned at startup, not based on real server capacity or load
* IP Hash uses plain modulo hashing — adding or removing a server reshuffles a large portion of client-to-server mappings, unlike consistent hashing
* `X-Forwarded-For` / `X-Real-IP` headers are trusted unconditionally; there's no allowlist of trusted upstream proxies, so a client could spoof its own routing key
* No automated test suite
* No TLS/HTTPS, graceful shutdown, or metrics/observability endpoints

## License

This project is for learning and experimentation.