package api

import (
	"Tiny-Proxy/internal/models"
	"Tiny-Proxy/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type APIServer struct {
	ProxyService *services.TCPProxyService
}

func NewAPIServer(proxyService *services.TCPProxyService) *APIServer {
	return &APIServer{
		ProxyService: proxyService,
	}
}

// StartAPIServer starts an HTTP server for managing domain mappings and proxy mappings
func (api *APIServer) Start(addr string) {
	http.HandleFunc("/domains", api.handleDomains)
	http.HandleFunc("/mappings", api.handleMappings)

	log.Printf("API server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("API server error: %v", err)
	}
}

// ExtractClientIP extracts the client's IP address from the HTTP request
func ExtractClientIP(r *http.Request) string {
	// Check if the request has the X-Forwarded-For header (common in proxies)
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0] // Use the first IP in the list
	}

	// Otherwise, use the remote address directly
	ip := r.RemoteAddr
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex] // Remove port
	}
	return ip
}

// handleDomains manages domain mapping CRUD operations
func (api *APIServer) handleDomains(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.getDomainMappings(w)
	case http.MethodPost:
		api.updateDomainMapping(w, r)
	case http.MethodDelete:
		api.deleteDomainMapping(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// handleMappings manages proxy mapping CRUD operations
func (api *APIServer) handleMappings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.getProxyMappings(w)
	case http.MethodPost:
		api.addProxyMapping(w, r)
	case http.MethodPut, http.MethodPatch:
		api.updateProxyMapping(w, r)
	case http.MethodDelete:
		api.deleteProxyMapping(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// getDomainMappings retrieves all domain mappings
func (api *APIServer) getDomainMappings(w http.ResponseWriter) {
	domainMap := api.ProxyService.GetDomainMappings()
	json.NewEncoder(w).Encode(domainMap)
}

// updateDomainMapping adds or updates a domain mapping
func (api *APIServer) updateDomainMapping(w http.ResponseWriter, r *http.Request) {
	var req models.DomainMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Replace "me" with the client's IP
	clientIP := ExtractClientIP(r)
	if req.IP == "me" {
		req.IP = clientIP
	}

	api.ProxyService.UpdateDomainMapping(req.Domain, req.IP)
	log.Printf("Updated domain mapping: %s -> %s", req.Domain, req.IP)
	w.WriteHeader(http.StatusNoContent)
}
// deleteDomainMapping removes a domain mapping
func (api *APIServer) deleteDomainMapping(w http.ResponseWriter, r *http.Request) {
	var req models.DomainMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := api.ProxyService.DeleteDomainMapping(req.Domain); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleted domain mapping: %s", req.Domain)
	w.WriteHeader(http.StatusNoContent)
}

// getProxyMappings retrieves all proxy mappings
func (api *APIServer) getProxyMappings(w http.ResponseWriter) {
	mappings := api.ProxyService.GetProxyMappings()
	json.NewEncoder(w).Encode(mappings)
}

// addProxyMapping adds a new proxy mapping
func (api *APIServer) addProxyMapping(w http.ResponseWriter, r *http.Request) {
	var req models.ProxyMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Replace "me" with the client's IP
	clientIP := ExtractClientIP(r)
	if strings.HasPrefix(req.ForwardAddr, "me:") {
		req.ForwardAddr = clientIP + req.ForwardAddr[2:] // Replace "me" with the IP
	}

	if err := api.ProxyService.AddMapping(req.ListenAddr, req.ForwardAddr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Added new mapping: %s -> %s", req.ListenAddr, req.ForwardAddr)
	w.WriteHeader(http.StatusNoContent)
}

// updateProxyMapping updates an existing proxy mapping's forward address
func (api *APIServer) updateProxyMapping(w http.ResponseWriter, r *http.Request) {
	var req models.ProxyMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Replace "me" with the client's IP
	clientIP := ExtractClientIP(r)
	if strings.HasPrefix(req.ForwardAddr, "me:") {
		req.ForwardAddr = clientIP + req.ForwardAddr[2:] // Replace "me" with the IP
	}

	if err := api.ProxyService.UpdateMapping(req.ListenAddr, req.ForwardAddr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updated mapping for listen address %s to forward to %s", req.ListenAddr, req.ForwardAddr)
	w.WriteHeader(http.StatusNoContent)
}
// deleteProxyMapping removes a proxy mapping
func (api *APIServer) deleteProxyMapping(w http.ResponseWriter, r *http.Request) {
	var req models.ProxyMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := api.ProxyService.DeleteMapping(req.ListenAddr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleted mapping for listen address: %s", req.ListenAddr)
	w.WriteHeader(http.StatusNoContent)
}