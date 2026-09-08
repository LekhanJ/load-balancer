package internal

import "sync"

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