package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL *url.URL
	IsHealthy bool
	weight int8
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