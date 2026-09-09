package algorithms

import (
	"sync"

	"github.com/LekhanJ/load-balancer/internal/balancer"
)

type RoundRobin struct {
	current int
	mu      sync.Mutex
}

func (r *RoundRobin) Next(servers []*balancer.Server) *balancer.Server {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 0; i < len(servers); i++ {
		idx := r.current % len(servers)
		next := servers[idx]
		r.current++

		if next.Healthy() {
			return next
		}
	}

	return nil
}