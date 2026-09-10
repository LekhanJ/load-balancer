package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/LekhanJ/load-balancer/internal/algorithms"
	"github.com/LekhanJ/load-balancer/internal/balancer"
)

func main() {
	var config balancer.Config
	lb := balancer.LoadBalancer{Strategy: &algorithms.LeastConnections{}}
	var servers []*balancer.Server
	set := make(map[int]struct{})

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	for _, serverURL := range config.Servers {
		parsedURL, err := url.Parse(serverURL)
		if err != nil {
			log.Fatal(err)
		}

		servers = append(servers, &balancer.Server{
			URL: parsedURL,
			Weight: int8(GetRandomWeight(&set, len(config.Servers))),
		})
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
		server := lb.GetNextServer(servers)

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