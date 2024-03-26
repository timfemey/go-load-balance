# Load Balancer

This is a load balancer well written for Go HTTP servers that supports various load balancing algorithms and supports concurrncy. It can be easily imported and used in the Go `main.go`.

## Usage

To use the load balancer in your `main.go`, follow these steps:

1. Import the load balancer package:

   ```go
   import (
       "fmt"
       "net/http"
       "net/url"
       "github.com/timfemey/go-load-balance/loadbalancer"
   )
   ```

2. Define backend servers as `[]string` objects:

   ```go
   endpoints := []string{
       "http://api1.example.com",
       "http://api2.example.com",

   }
   ```

3. Create an instance of the load balancer with the desired load balancing algorithm:

   ```go
   lb := loadbalancer.ResourceBasedLoadBalancer(endpoints)
   // OR
   lb := loadbalancer.LeastConnectionsLoadBalancer(endpoints)
   ```

4. Start the HTTP server with the load balancer:

   ```go
   fmt.Println("Load balancer listening on :8080")
   http.ListenAndServe(":8080", lb)
   ```

## Examples

### Resource-Based Load Balancing

```go
endpoints := []string{
   "http://api1.example.com",
   "http://api2.example.com",
    // Add more backend servers as needed
}

lb := loadbalancer.NewResourceBasedLoadBalancer(endpoints)
fmt.Println("Load balancer listening on :8080")
http.ListenAndServe(":8080", lb)
```

### Least Connections Load Balancing

```go
endpoints := []string{
    "http://api1.example.com",
    "http://api2.example.com",

}

lb := loadbalancer.NewLeastConnectionsLoadBalancer(endpoints)
fmt.Println("Load balancer listening on :8080")
http.ListenAndServe(":8080", lb)
```

## Load Balancing Algorithms

### Resource-Based Load Balancing

This algorithm periodically hits the /health endpoint of each server to check CPU and memory resources. Requests are directed to the server with the best resource utilization.

### Least Connections Load Balancing

This algorithm keeps track of the number of active connections to each server and directs new requests to the server with the fewest active connections.

### Weighted Round-Robin Load Balancing

The weighted round-robin load balancing algorithm assigns weights to each backend server and distributes requests in a round-robin manner, with each server receiving requests proportionally to its weight.

### Sticky Round-Robin Load Balancing

The sticky round-robin load balancing algorithm ensures that requests from the same client are consistently directed to the same backend server, providing session affinity for stateful applications.

### Hash-Based Load Balancing

The hash-based load balancing algorithm maps incoming requests to backend servers based on a hash function applied to a characteristic of the request, such as the client's IP address or the request URL.

## Contributing

Feel free to contribute to the development of this HTTP Load Balancer. Create issues, submit pull requests, and help make it even better!

## Acknowledgments

- Thanks to the Golang community for providing a powerful language for building efficient and concurrent applications.
