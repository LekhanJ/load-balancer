package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/LekhanJ/load-balancer/internal/algorithms"
	"github.com/LekhanJ/load-balancer/internal/balancer"
)

func loadConfig(path string) (*balancer.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config balancer.Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func selectAlgorithm(name string) balancer.Algorithm {
	switch name {
	case "round-robin":
		return &algorithms.RoundRobin{}
	case "weighted-round-robin":
		return &algorithms.WeightedRoundRobin{}
	case "least-connections":
		return &algorithms.LeastConnections{}
	case "ip-hash":
		return &algorithms.IPHash{}
	default:
		log.Fatalf("unknown algorithm: %s", name)
		return nil
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}


func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	lb := balancer.LoadBalancer{Strategy: selectAlgorithm(config.Algorithm)}

	var servers []*balancer.Server
	set := make(map[int]struct{})
	for _, serverURL := range config.Servers {
		parsedURL, err := url.Parse(serverURL)
		if err != nil {
			log.Fatal(err)
		}

		server := &balancer.Server{
			URL: parsedURL,
			Weight: int8(GetRandomWeight(&set, len(config.Servers))),
		}
		server.SetHealthy(true)
		servers = append(servers, server)
	}

	interval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		go balancer.HealthCheck(server, interval)
	}

	pool := balancer.NewProxyPool()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		server := lb.GetNextServer(servers, key)

		if server == nil {
			http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
			return
		}

		server.IncrementConnections()
    	defer server.DecrementConnections()

		w.Header().Add("X-Forwarded-Server", server.URL.String())
		pool.Get(server).ServeHTTP(w, r)
	})

	log.Println("Starting load balancer on port", config.Port)

	err = http.ListenAndServe(config.Port, nil)
	if err != nil {
			log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}
}

func GetRandomWeight(set *map[int]struct{}, limit int) int {
	for {
		num := rand.IntN(limit) + 1
		if _, exists := (*set)[num]; !exists {
			(*set)[num] = struct{}{}
			return num
		}
	}
}