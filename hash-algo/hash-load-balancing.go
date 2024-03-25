package main

import (
	"hash/fnv"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type hashLoadBalancingStruct struct {
	endpoints []string
	mu        sync.Mutex
}

func HashLoadBalancing(endpoints []string) *hashLoadBalancingStruct {
	return &hashLoadBalancingStruct{
		endpoints: endpoints,
	}
}

func (hlb *hashLoadBalancingStruct) hashFunc(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))

	return h.Sum32()
}

func (hlb *hashLoadBalancingStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hlb.mu.Lock()
	defer hlb.mu.Unlock()

	clientIP := r.RemoteAddr

	hash := hlb.hashFunc(clientIP)

	index := int(hash) % len(hlb.endpoints)

	proxyURL, _ := url.Parse(hlb.endpoints[index])
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}
