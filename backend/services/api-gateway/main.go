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
		userServiceURL:     "http://user-service:8001",
		paymentServiceURL:  "http://payment-service:8002",
		strategyServiceURL: "http://strategy-service:8003",
	}
}

// healthHandler - Gateway's own health check
func (g *Gateway) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"service": "API Gateway - Train Service",
		"status":  "healthy",
		"port":    8080,
		"routes": map[string]string{
			"/user/*":      "forwards to Customer Service (port 8001)",
			"/users/*":     "forwards to Customer Service (port 8001)",
			"/payments/*":  "forwards to Payment Service (port 8002)",
			"/strategies/*": "forwards to Strategy Service (port 8003)",
		},
	}
	json.NewEncoder(w).Encode(response)
}

// welcomeHandler - Gateway welcome message
func (g *Gateway) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "🚂 Welcome to the Train Service (API Gateway)!",
		"description": "I route customers to different buildings in your shopping complex",
		"services": map[string]string{
			"Customer Service": "http://localhost:8080/user/ OR http://localhost:8080/users/",
			"Payment Service":  "http://localhost:8080/payments/",
			"Strategy Service": "http://localhost:8080/strategies/",
		},
		"health_check": "http://localhost:8080/health",
	}
	json.NewEncoder(w).Encode(response)
}

// createProxy creates a reverse proxy to forward requests
func (g *Gateway) createProxy(targetURL string) *httputil.ReverseProxy {
	target, _ := url.Parse(targetURL)
	return httputil.NewSingleHostReverseProxy(target)
}

// proxyHandler routes requests to appropriate microservice
func (g *Gateway) proxyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	var targetService string
	var newPath string
	
	switch {
	case strings.HasPrefix(path, "/users/") || strings.HasPrefix(path, "/user/"):
		targetService = g.userServiceURL
		if strings.HasPrefix(path, "/users/") {
			newPath = strings.TrimPrefix(path, "/users")
		} else {
			newPath = strings.TrimPrefix(path, "/user")
		}
		
	case strings.HasPrefix(path, "/payments/"):
		targetService = g.paymentServiceURL
		newPath = strings.TrimPrefix(path, "/payments")
		
	case strings.HasPrefix(path, "/strategies/"):
		targetService = g.strategyServiceURL
		newPath = strings.TrimPrefix(path, "/strategies")
		
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Service not found. Available routes: /user/, /users/, /payments/, /strategies/",
		})
		return
	}
	
	r.URL.Path = newPath
	if newPath == "" {
		r.URL.Path = "/"
	}
	
	proxy := g.createProxy(targetService)
	log.Printf("🚂 Routing %s %s -> %s%s", r.Method, path, targetService, r.URL.Path)
	proxy.ServeHTTP(w, r)
}

func main() {
	gateway := NewGateway()
	
	// Single handler for all routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// DEBUG: Log every request
		log.Printf("📝 Request: %s %s", r.Method, path)
		
		switch {
		case path == "/":
			log.Printf("📝 Routing to welcomeHandler")
			gateway.welcomeHandler(w, r)
		case path == "/health":
			log.Printf("📝 Routing to healthHandler")
			gateway.healthHandler(w, r)
		case strings.HasPrefix(path, "/users/") || strings.HasPrefix(path, "/user/"):
			log.Printf("📝 Routing to proxyHandler for user service")
			gateway.proxyHandler(w, r)
		case strings.HasPrefix(path, "/payments/"):
			log.Printf("📝 Routing to proxyHandler for payment service")
			gateway.proxyHandler(w, r)
		case strings.HasPrefix(path, "/strategies/"):
			log.Printf("📝 Routing to proxyHandler for strategy service")
			gateway.proxyHandler(w, r)
		default:
			log.Printf("📝 Unknown route: %s", path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Service not found. Available routes: /user/, /users/, /payments/, /strategies/",
			})
		}
	})
	
	port := ":8080"
	fmt.Printf("🚂 NEW VERSION - Train Service (API Gateway) starting on port%s\n", port)
	fmt.Println("📋 Available routes:")
	fmt.Println("   GET  / - Gateway welcome message")
	fmt.Println("   GET  /health - Gateway health check")
	fmt.Println("   ALL  /user/* - Forward to Customer Service (port 8001)")
	fmt.Println("   ALL  /users/* - Forward to Customer Service (port 8001)")
	fmt.Println("   ALL  /payments/* - Forward to Payment Service (port 8002)")
	fmt.Println("   ALL  /strategies/* - Forward to Strategy Service (port 8003)")
	fmt.Println()
	fmt.Printf("🌐 Try: http://localhost%s\n", port)
	
	log.Fatal(http.ListenAndServe(port, nil))
}
