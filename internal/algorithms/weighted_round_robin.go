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
	wr.mu.Lock()
	defer wr.mu.Unlock()

	var totalWeight int16
	var weightedServer *balancer.Server

	for i:=0; i < len(servers); i++ {
		if servers[i].Healthy() {
			totalWeight += int16(servers[i].Weight)
		}
	}

	for i:=0; i < len(servers); i++ {
		if servers[i].Healthy() {
			servers[i].CurrWeight += servers[i].Weight
			if weightedServer == nil || servers[i].CurrWeight > weightedServer.CurrWeight {
				weightedServer = servers[i]
			}
		}
	}

	weightedServer.CurrWeight -= int8(totalWeight)

	return weightedServer
}