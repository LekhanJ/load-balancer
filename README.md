# Load Balancer in Go

A simple HTTP load balancer built with **Go**, using **Python backend servers**.

This project is based on the article [Building a Simple Load Balancer in Go](https://dev.to/vivekalhat/building-a-simple-load-balancer-in-go-70d).

## Overview

The load balancer receives HTTP requests and forwards them to healthy backend servers. Load-balancing strategies are implemented behind a common `Algorithm` interface, so the distribution logic can be swapped out independently of the proxying and health-check machinery. The binary currently starts up with **Least Connections** as its active strategy.

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
internal/balancer/     Server model, config struct, health checks, proxy pooling
internal/algorithms/   Load-balancing strategies (Algorithm implementations)
server/server.py       Minimal Python backend used for local testing
```

## Features

* Pluggable load-balancing strategies via a shared `Algorithm` interface
* **Least Connections** — routes to the healthy server with the fewest active connections (used by default in `main.go`)
* **Round Robin** — cycles through healthy servers in order
* **Weighted Round Robin** — favors servers with a higher assigned weight
* **IP Hash** — scaffolded (`internal/algorithms/ip_hash.go`) but not yet implemented
* Per-server active-connection tracking, used by the Least Connections strategy
* Periodic backend health checks, with automatic skipping of unhealthy servers
* Reverse-proxy connection pooling (`ProxyPool`) so each server reuses a single `httputil.ReverseProxy`
* Simple configuration through JSON

## Configuration

The backend servers and load balancer port are configured in `cmd/load-balancer/config.json`:

```json
{
  "port": ":8081",
  "healthCheckInterval": "2s",
  "servers": [
    "http://localhost:5001",
    "http://localhost:5002",
    "http://localhost:5003",
    "http://localhost:5004",
    "http://localhost:5005"
  ]
}
```

Note: the active strategy (currently `LeastConnections`) is set in `cmd/load-balancer/main.go`, not in `config.json` — see [Future Improvements](#future-improvements).

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

1. The load balancer reads the backend server configuration and assigns each server a randomized weight (used by the weighted strategy).
2. It starts a goroutine per server that periodically checks whether that backend is healthy.
3. For every incoming request, the active `Algorithm` selects the next healthy server.
4. The request is forwarded through a pooled `httputil.ReverseProxy` for that server, with active-connection counts incremented/decremented around the call.
5. If no healthy servers are available, it returns `503 Service Unavailable`.

## Learning Goals

This project is intended to understand:

* HTTP request forwarding
* Reverse proxies
* Load-balancing strategies (round-robin, weighted round-robin, least connections, IP hash)
* Health checks
* Goroutines and concurrency
* Mutexes and shared state
* Basic Go networking

## License

This project is for learning and experimentation.