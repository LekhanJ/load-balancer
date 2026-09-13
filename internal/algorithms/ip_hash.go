package algorithms

import (
	"hash/fnv"

	"github.com/LekhanJ/load-balancer/internal/balancer"
)

type IPHash struct{}

func hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (ih *IPHash) Next(servers []*balancer.Server, key string) *balancer.Server {
	healthy := make([]*balancer.Server, 0, len(servers))
	for _, s := range servers {
		if s.Healthy() {
			healthy = append(healthy, s)
		}
	}
	if len(healthy) == 0 {
		return nil
	}

	idx := hashKey(key) % uint32(len(healthy))
	return healthy[idx]
}