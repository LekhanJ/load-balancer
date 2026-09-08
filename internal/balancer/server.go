package internal

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL *url.URL
	IsHealthy bool
	mu sync.Mutex
}

func (s *Server) ReverseProxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}