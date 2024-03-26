package loadbalancer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type ServerResponse struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

type resourceBasedLoadBalancer struct {
	endpoints  []string
	bestServer string
	mu         sync.Mutex
}

func (rlb *resourceBasedLoadBalancer) updateBestServer() {
	for {
		var bestServerURL string
		var bestScore float64 = -1

		for _, endpoint := range rlb.endpoints {
			resource, err := rlb.getServerResource(endpoint)
			if err != nil {
				fmt.Printf("Failed to get resource for server %s: %v\n", endpoint, err)
				continue
			}

			score := rlb.calculateScore(resource)
			if score > bestScore {
				bestScore = score
				bestServerURL = endpoint
			}
		}

		rlb.mu.Lock()
		rlb.bestServer = bestServerURL
		rlb.mu.Unlock()

		time.Sleep(30 * time.Second)
	}
}

func (rlb *resourceBasedLoadBalancer) getServerResource(endpoint string) (*ServerResponse, error) {
	resp, err := http.Get(endpoint + "/health")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Health Check Failed, Server returned non-200 staus code ", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resource ServerResponse
	err = json.Unmarshal(body, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (rlb *resourceBasedLoadBalancer) calculateScore(resource *ServerResponse) float64 {
	return resource.CPUUsage + resource.MemoryUsage
}

func (rlb *resourceBasedLoadBalancer) getBestServerURL() string {
	rlb.mu.Lock()
	defer rlb.mu.Unlock()
	return rlb.bestServer
}

func (rlb *resourceBasedLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rlb.mu.Lock()
	defer rlb.mu.Unlock()

	if rlb.bestServer == "" {
		rlb.updateBestServer()
	}

	proxyURL, _ := url.Parse(rlb.getBestServerURL())
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
}

func ResourceBasedLoadBalancer(endpoints []string) *resourceBasedLoadBalancer {
	return &resourceBasedLoadBalancer{
		endpoints: endpoints,
	}
}
