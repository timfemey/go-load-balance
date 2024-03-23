package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type roundRobinLoadBalancerStruct struct {
	endpoints []string
	index     int
	mu        sync.Mutex
}

func RoundRobinLoadBalancer(endpoints []string) *roundRobinLoadBalancerStruct {
	return &roundRobinLoadBalancerStruct{
		endpoints: endpoints,
		index:     0,
	}
}

func (lb *roundRobinLoadBalancerStruct) ServeHTTP(w http.ResponseWriter, r http.Request) {
	lb.mu.Lock()
	endpoint := lb.endpoints[lb.index]
	lb.index = (lb.index + 1) % len(lb.endpoints)
	lb.mu.Unlock()

	proxyURL, _ := url.Parse(endpoint)
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, &r)
}
