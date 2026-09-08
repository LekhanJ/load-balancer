package main

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"

	"github.com/LekhanJ/load-balancer/internal/balancer"
)

func main() {
	var config balancer.Config
	var lb balancer.LoadBalancer
	var servers []*balancer.Server

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