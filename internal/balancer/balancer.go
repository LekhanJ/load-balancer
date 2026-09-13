package balancer

type Algorithm interface {
	Next(servers []*Server, key string) *Server
}
 
type LoadBalancer struct {
	Strategy Algorithm
}
 
func (lb *LoadBalancer) GetNextServer(servers []*Server, key string) *Server {
	return lb.Strategy.Next(servers, key)
}