# Load Balancer in Go

A simple HTTP load balancer built with **Go**, using **Python backend servers**.

This project is based on the article [Building a Simple Load Balancer in Go](https://dev.to/vivekalhat/building-a-simple-load-balancer-in-go-70d).

## Overview

The load balancer receives HTTP requests and forwards them to healthy backend servers using a **round-robin** strategy.

```text
Client
   |
   v
Go Load Balancer (:8080)
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

## Features

* Round-robin request distribution
* Periodic backend health checks
* Automatic skipping of unhealthy servers
* HTTP reverse proxy using Go's standard library
* Simple configuration through JSON

## Configuration

The backend servers and load balancer port are configured in `config.json`.

```json
{
  "port": ":8080",
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

## Running the Project

### 1. Start the Python backend servers

Open five terminal windows and run:

```bash
PORT=5001 python3 backend/server.py
PORT=5002 python3 backend/server.py
PORT=5003 python3 backend/server.py
PORT=5004 python3 backend/server.py
PORT=5005 python3 backend/server.py
```

Each server listens on its own port.

### 2. Start the Go load balancer

From the project root:

```bash
go run .
```

The load balancer starts on:

```text
http://localhost:8080
```

### 3. Send requests

```bash
curl http://localhost:8080
```

Example response:

```text
Response from Python server 5001
```

Send multiple requests to observe round-robin distribution:

```bash
for i in {1..10}; do curl http://localhost:8080; done
```

## How It Works

1. The load balancer reads the backend server configuration.
2. It periodically checks whether each backend is healthy.
3. For every incoming request, it selects the next healthy server.
4. The request is forwarded using `httputil.ReverseProxy`.
5. If no healthy servers are available, it returns `503 Service Unavailable`.

## Learning Goals

This project is intended to understand:

* HTTP request forwarding
* Reverse proxies
* Round-robin load balancing
* Health checks
* Goroutines and concurrency
* Mutexes and shared state
* Basic Go networking

## Future Improvements

* Configurable load-balancing strategies
* Better health-check handling
* Request timeouts
* Graceful shutdown
* Metrics and logging
* Dynamic backend registration
* Docker support

## License

This project is for learning and experimentation.
