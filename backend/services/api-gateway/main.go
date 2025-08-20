package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Gateway represents our Train Service
type Gateway struct {
	userServiceURL     string
	paymentServiceURL  string
	strategyServiceURL string
}

// NewGateway creates a new API Gateway instance
func NewGateway() *Gateway {
	return &Gateway{
		userServiceURL:     "http://user-service-container:8001",
		paymentServiceURL:  "http://payment-service-container:8002",
		strategyServiceURL: "http://strategy-service-container:8003",
	}
}

// healthHandler - Gateway's own health check
func (g *Gateway) healthHandler(w http.ResponseWriter, r *http.Request) {
	// Set response headers (like setting the response type)
	w.Header().Set("Content-Type", "application/json")
	
	// Create response data (like Python dictionary)
	response := map[string]interface{}{
		"service": "API Gateway - Train Service",
		"status":  "healthy",
		"port":    8080,
		"routes": map[string]string{
			"/users/*":     "forwards to Customer Service (port 8001)",
			"/payments/*":  "forwards to Payment Service (port 8002)",
			"/strategies/*": "forwards to Strategy Service (port 8003)",
		},
	}
	
	// Convert to JSON and send (like FastAPI's automatic return)
	json.NewEncoder(w).Encode(response)
}

// welcomeHandler - Gateway welcome message
func (g *Gateway) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := map[string]interface{}{
		"message": "🚂 Welcome to the Train Service (API Gateway)!",
		"description": "I route customers to different buildings in your shopping complex",
		"services": map[string]string{
			"Customer Service": "http://localhost:8080/users/",
			"Payment Service":  "http://localhost:8080/payments/",
			"Strategy Service": "http://localhost:8080/strategies/",
		},
		"health_check": "http://localhost:8080/health",
	}
	
	json.NewEncoder(w).Encode(response)
}

// createProxy creates a reverse proxy to forward requests
func (g *Gateway) createProxy(targetURL string) *httputil.ReverseProxy {
	// Parse the target service URL
	target, _ := url.Parse(targetURL)
	
	// Create reverse proxy (this forwards requests)
	return httputil.NewSingleHostReverseProxy(target)
}

// proxyHandler routes requests to appropriate microservice
func (g *Gateway) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the request path
	path := r.URL.Path
	
	// Determine which service to forward to (like train routing)
	var targetService string
	var newPath string
	
	switch {
	case strings.HasPrefix(path, "/users/"):
		targetService = g.userServiceURL
		newPath = strings.TrimPrefix(path, "/users")
		
	case strings.HasPrefix(path, "/payments/"):
		targetService = g.paymentServiceURL
		newPath = strings.TrimPrefix(path, "/payments")
		
	case strings.HasPrefix(path, "/strategies/"):
		targetService = g.strategyServiceURL
		newPath = strings.TrimPrefix(path, "/strategies")
		
	default:
		// Unknown route
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Service not found. Available routes: /users/, /payments/, /strategies/",
		})
		return
	}
	
	// Update the request path for the target service
	r.URL.Path = newPath
	if newPath == "" {
		r.URL.Path = "/"
	}
	
	// Create proxy and forward the request
	proxy := g.createProxy(targetService)
	
	// Add logging (like seeing the train route)
	log.Printf("🚂 Routing %s %s -> %s%s", r.Method, path, targetService, r.URL.Path)
	
	// Forward the request
	proxy.ServeHTTP(w, r)
}

func main() {
	// Create our Train Service instance
	gateway := NewGateway()
	
	// Set up routes (like train schedule)
	http.HandleFunc("/", gateway.welcomeHandler)           // Main station info
	http.HandleFunc("/health", gateway.healthHandler)     // Train service status
	http.HandleFunc("/users/", gateway.proxyHandler)      // Route to Customer Service
	http.HandleFunc("/payments/", gateway.proxyHandler)   // Route to Payment Service
	http.HandleFunc("/strategies/", gateway.proxyHandler) // Route to Strategy Service
	
	// Start the Train Service
	port := ":8080"
	fmt.Printf("🚂 Train Service (API Gateway) starting on port%s\n", port)
	fmt.Println("📋 Available routes:")
	fmt.Println("   GET  / - Gateway welcome message")
	fmt.Println("   GET  /health - Gateway health check")
	fmt.Println("   ALL  /users/* - Forward to Customer Service (port 8001)")
	fmt.Println("   ALL  /payments/* - Forward to Payment Service (port 8002)")
	fmt.Println("   ALL  /strategies/* - Forward to Strategy Service (port 8003)")
	fmt.Println()
	fmt.Printf("🌐 Try: http://localhost%s\n", port)
	
	// Start HTTP server (like opening the train station)
	log.Fatal(http.ListenAndServe(port, nil))
}
