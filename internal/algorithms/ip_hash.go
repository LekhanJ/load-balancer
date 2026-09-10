package algorithms

import "github.com/LekhanJ/load-balancer/internal/balancer"

type IPHash struct{}

func (i *IPHash) Next(servers []*balancer.Server) *balancer.Server {
	
}