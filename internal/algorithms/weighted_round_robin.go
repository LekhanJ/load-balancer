package algorithms

import (
	"sync"

	"github.com/LekhanJ/load-balancer/internal/balancer"
)

type WeightedRoundRobin struct {
	current int
	mu      sync.Mutex
}

func (wr *WeightedRoundRobin) Next(servers []*balancer.Server) *balancer.Server {
	
}