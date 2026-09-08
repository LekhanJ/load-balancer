package internal

import (
	"fmt"
	"net/http"
	"time"
)

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