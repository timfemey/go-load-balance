package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type leastConnectionsLoadBalancer struct {
	endpoints         []string
	activeConnections map[string]int
	mu                sync.Mutex
}

func LeastConnectionsLoadBalancer(endpoints []string) *leastConnectionsLoadBalancer {
	activeConns := make(map[string]int)
	for _, endpoint := range endpoints {
		activeConns[endpoint] = 0
	}

	return &leastConnectionsLoadBalancer{
		endpoints:         endpoints,
		activeConnections: activeConns,
	}
}

func (lcb *leastConnectionsLoadBalancer) IncrementConnection(endpoint string) {
	lcb.mu.Lock()
	defer lcb.mu.Unlock()
	lcb.activeConnections[endpoint] += 1
}

func (lcb *leastConnectionsLoadBalancer) DecrementConnection(endpoint string) {
	lcb.mu.Lock()
	defer lcb.mu.Unlock()
	lcb.activeConnections[endpoint] -= 1
}

func (lcb *leastConnectionsLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var leastConns string
	minConns := int(^uint(0) >> 1)
	for _, endpoint := range lcb.endpoints {
		conns := lcb.activeConnections[endpoint]
		if minConns > conns {
			minConns = conns
			leastConns = endpoint
		}

	}
	proxyURL, _ := url.Parse(leastConns)
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}
