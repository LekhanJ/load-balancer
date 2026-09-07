package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "sync"
    "time"
)

type LoadBalancer struct {
	Current int
	mu sync.Mutex
}

func (lb *LoadBalancer) getNextServer(servers []*Server) *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i:=0; i<len(servers); i++ {
		idx := lb.Current % len(servers)
		nextServer := servers[idx]
		lb.Current++

		nextServer.mu.Lock()
		isHealthy := nextServer.IsHealthy
		nextServer.mu.Unlock()

		if isHealthy {
			return nextServer
		}
	}

	return nil
}

type Server struct {
	URL *url.URL
	IsHealthy bool
	mu sync.Mutex
}

func (s *Server) ReverseProxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}

type Config struct {
	Port string `json:"port"`
	HealthCheckInterval string `json:"healthCheckInterval"`
	Servers []string `json:"servers"`
}

func healthCheck(s *Server, interval time.Duration) {
    for range time.Tick(interval) {
        res, err := http.Head(s.URL.String())

        s.mu.Lock()
        if err != nil {
            s.IsHealthy = false
            fmt.Printf("%s is down\n", s.URL)
        } else {
            s.IsHealthy = res.StatusCode == http.StatusOK

            if !s.IsHealthy {
                fmt.Printf("%s is down\n", s.URL)
            }
        }
        s.mu.Unlock()
    }
}

func main() {
	var config Config
	var lb LoadBalancer
	var servers []*Server

	file, err := os.Open("../../config.json")
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

		servers = append(servers, &Server{
			URL: parsedURL,
		})
	}

	interval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		go healthCheck(server, interval)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := lb.getNextServer(servers)

		if server == nil {
			http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
			return
		}

		w.Header().Add("X-Forwarded-Server", server.URL.String())
		server.ReverseProxy().ServeHTTP(w, r)
	})

	log.Println("Starting load balancer on port", config.Port)

	err = http.ListenAndServe(config.Port, nil)
	if err != nil {
			log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}
}