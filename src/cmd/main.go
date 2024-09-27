package main

import (
	"Tiny-Proxy/internal/api"
	"Tiny-Proxy/internal/services"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <api_address:port>", os.Args[0])
	}

	apiAddr := os.Args[1]

	// Initialize the TCPProxy
	proxyServer := services.NewTCPProxy()

	// Create and start API server for managing domain and proxy mappings
	apiServer := api.NewAPIServer(proxyServer)
	apiServer.Start(apiAddr)
}