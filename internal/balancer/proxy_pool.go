package balancer

import (
	"net/http/httputil"
	"sync"
)

type ProxyPool struct {
    mu     sync.Mutex
    proxies map[*Server]*httputil.ReverseProxy
}

func NewProxyPool() *ProxyPool {
    return &ProxyPool{proxies: make(map[*Server]*httputil.ReverseProxy)}
}

func (p *ProxyPool) Get(s *Server) *httputil.ReverseProxy {
    p.mu.Lock()
    defer p.mu.Unlock()

    if proxy, ok := p.proxies[s]; ok {
        return proxy
    }
    proxy := httputil.NewSingleHostReverseProxy(s.URL)
    p.proxies[s] = proxy
    return proxy
}