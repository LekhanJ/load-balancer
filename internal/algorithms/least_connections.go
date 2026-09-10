package algorithms

import "github.com/LekhanJ/load-balancer/internal/balancer"

type LeastConnections struct{}

func (l *LeastConnections) Next(servers []*balancer.Server) *balancer.Server {
	var chosen *balancer.Server

    for _, s := range servers {
        if !s.Healthy() {
            continue
        }
        if chosen == nil || s.ActiveConnections() < chosen.ActiveConnections() {
            chosen = s
        }
    }

    return chosen
}