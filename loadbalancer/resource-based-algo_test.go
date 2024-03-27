package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestResourceBasedLoadBalancer(t *testing.T) {
	mockServers := []*httptest.Server{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Simulate server processing time
			w.WriteHeader(http.StatusOK)
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond) // Simulate server processing time
			w.WriteHeader(http.StatusOK)
		})),
	}

	// Create a resource-based load balancer
	lb := ResourceBasedLoadBalancer([]string{
		mockServers[0].URL,
		mockServers[1].URL,
	})

	// Send 100 requests and measure response time
	start := time.Now()
	for i := 0; i < 100; i++ {
		req, err := http.NewRequest("GET", "https://catfact.ninja/fact", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		w := httptest.NewRecorder()
		lb.ServeHTTP(w, req)
	}
	elapsed := time.Since(start)

	fmt.Printf("Resource Based Load balancing: Average response time = %v\n", elapsed/100)

	// Close mock servers
	for _, server := range mockServers {
		server.Close()
	}
}
