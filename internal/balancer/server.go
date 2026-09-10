package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL *url.URL
	IsHealthy bool
	Weight int8
	CurrWeight int8
	activeConnections int
	mu sync.Mutex
}

func (s *Server) ReverseProxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}

func (s *Server) Healthy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.IsHealthy
}

func (s *Server) SetHealthy(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsHealthy = healthy
}

func (s *Server) ActiveConnections() int {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.activeConnections
}

func (s *Server) IncrementConnections() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.activeConnections++
}

func (s *Server) DecrementConnections() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.activeConnections--
}