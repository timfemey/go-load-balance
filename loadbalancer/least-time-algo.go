package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type leastTimeLoadBalancer struct {
	endpoints      []string
	responsesTimes map[string]time.Duration
	mu             sync.Mutex
}

func LeastTimeLoadBalancer(endpoints []string) *leastTimeLoadBalancer {
	responsesTimes := make(map[string]time.Duration)
	for _, endpoint := range endpoints {
		responsesTimes[endpoint] = 0
	}

	return &leastTimeLoadBalancer{
		endpoints:      endpoints,
		responsesTimes: responsesTimes,
	}
}

func (ltb *leastTimeLoadBalancer) UpdateResponseTime(endpoint string, responseTime time.Duration) {
	ltb.mu.Lock()
	defer ltb.mu.Unlock()

	ltb.responsesTimes[endpoint] = responseTime
}

func (ltb *leastTimeLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var leastTimeServer string
	minTime := time.Duration(^uint64(0) >> 1)

	for _, endpoint := range ltb.endpoints {
		ltb.mu.Lock()
		responseTime := ltb.responsesTimes[endpoint]
		ltb.mu.Unlock()

		if minTime > responseTime {
			minTime = responseTime
			leastTimeServer = endpoint
		}
	}

	proxyURL, _ := url.Parse(leastTimeServer)
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}
