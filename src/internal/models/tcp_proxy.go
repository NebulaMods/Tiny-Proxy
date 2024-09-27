package models

import (
	"net"
	"sync"
)

// TCPProxy represents the TCP proxy manager with multiple mappings
type TCPProxy struct {
	DomainMap map[string]string     // Custom domain to IP mappings
	Mappings  map[string]ProxyMapping // Mapping from listener to forward address
	Listeners map[string]net.Listener // Store active listeners
	Mu        sync.RWMutex            // Mutex for thread-safe access
}
