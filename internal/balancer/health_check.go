package balancer

import (
	"fmt"
	"net/http"
	"time"
)

func HealthCheck(s *Server, interval time.Duration) {
    for range time.Tick(interval) {
        res, err := http.Head(s.URL.String())
 
        if err != nil {
            s.SetHealthy(false)
            fmt.Printf("%s is down\n", s.URL)
        } else {
            healthy := res.StatusCode == http.StatusOK
            s.SetHealthy(healthy)
 
            if !healthy {
                fmt.Printf("%s is down\n", s.URL)
            }
        }
    }
}