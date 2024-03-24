package main

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type stickyROundRobinLoadBalancerStruct struct {
	endpoints []string          // Endpoints
	index     map[string]string //Map Clinet IP to endpoint
	mu        sync.Mutex
}

func StickyROundRobinLoadBalancer(endpoints []string) *stickyROundRobinLoadBalancerStruct {
	return &stickyROundRobinLoadBalancerStruct{
		endpoints: endpoints,
		index:     make(map[string]string),
	}
}

func (lb *stickyROundRobinLoadBalancerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	lb.mu.Lock()
	defer lb.mu.Unlock()

	endpoint, ok := lb.index[clientIP]
	if !ok {
		max := len(lb.endpoints) - 1
		min := 0
		randomNum := rand.Intn(max-min) + min
		endpoint = lb.endpoints[randomNum]
		lb.index[clientIP] = endpoint
	}

	proxyURL, _ := url.Parse(endpoint)
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}
