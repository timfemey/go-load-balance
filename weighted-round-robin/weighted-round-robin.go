package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type weightedEndpoint struct {
	endpoint string // A valid URL Endpoint
	weight   int    // A weight for each endpoint
}

type weightedRoundRobinLoadBalancerStruct struct {
	endpoints   []weightedEndpoint
	totalWeight int
	curr        int
	mu          sync.Mutex
}

func WeightedRoundRobinLoadBalancer(endpoints []weightedEndpoint) *weightedRoundRobinLoadBalancerStruct {
	lb := &weightedRoundRobinLoadBalancerStruct{
		endpoints: endpoints,
	}
	lb.calculateTotalWeight()
	return lb
}

func (lb *weightedRoundRobinLoadBalancerStruct) calculateTotalWeight() {
	lb.totalWeight = 0
	tempCarrier := 0
	for _, endpoint := range lb.endpoints {
		tempCarrier += endpoint.weight
	}
	lb.totalWeight = tempCarrier
}

func (lb *weightedRoundRobinLoadBalancerStruct) selectNextIndex() int {
	lb.curr = (lb.curr + 1) % lb.totalWeight
	currWeight := lb.curr
	for i, endpoint := range lb.endpoints {
		currWeight -= endpoint.weight
		if 0 > currWeight {
			return i
		}
	}

	return 0
}

func (lb *weightedRoundRobinLoadBalancerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	index := lb.selectNextIndex()
	endpoint := lb.endpoints[index].endpoint

	proxyURL, _ := url.Parse(endpoint)
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}
