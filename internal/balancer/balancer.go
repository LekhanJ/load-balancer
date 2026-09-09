package balancer

type Algorithm interface {
	Next(servers []*Server) *Server
}
 
type LoadBalancer struct {
	Strategy Algorithm
}
 
func (lb *LoadBalancer) GetNextServer(servers []*Server) *Server {
	return lb.Strategy.Next(servers)
}