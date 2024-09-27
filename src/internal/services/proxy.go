package services

import (
	"Tiny-Proxy/internal/models"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// TCPProxyService extends the models.TCPProxy to add functionalities
type TCPProxyService struct {
	*models.TCPProxy
}

// NewTCPProxy initializes a new TCPProxyService instance
func NewTCPProxy() *TCPProxyService {
	return &TCPProxyService{
		TCPProxy: &models.TCPProxy{
			DomainMap: make(map[string]string),
			Mappings:  make(map[string]models.ProxyMapping),
			Listeners: make(map[string]net.Listener),
		},
	}
}

// AddMapping adds a new mapping and starts listening on the specified port
func (p *TCPProxyService) AddMapping(listenAddr, forwardAddr string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if _, exists := p.Mappings[listenAddr]; exists {
		return errors.New("mapping already exists for this listen address")
	}

	// Start a listener on the given address
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	// Store mapping and listener
	p.Mappings[listenAddr] = models.ProxyMapping{ListenAddr: listenAddr, ForwardAddr: forwardAddr}
	p.Listeners[listenAddr] = listener

	go p.startListener(listener, forwardAddr)

	log.Printf("Listening on %s, forwarding to %s\n", listenAddr, forwardAddr)
	return nil
}

// startListener accepts incoming connections and forwards them to the forward address
func (p *TCPProxyService) startListener(listener net.Listener, forwardAddr string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go p.handleConnection(conn, forwardAddr)
	}
}

// handleConnection handles the client connection and forwards traffic
func (p *TCPProxyService) handleConnection(srcConn net.Conn, forwardAddr string) {
	defer srcConn.Close()

	// Resolve the forward address to an IP address
	resolvedAddr, err := p.resolveForwardAddress(forwardAddr)
	if err != nil {
		log.Println("Error resolving address:", err)
		return
	}

	// Establish a connection to the destination
	dstConn, err := net.DialTimeout("tcp", resolvedAddr.String(), 30*time.Second)
	if err != nil {
		log.Println("Error connecting to forward address:", err)
		return
	}
	defer dstConn.Close()

	// Use a WaitGroup to ensure both directions are copied before closing
	var wg sync.WaitGroup
	wg.Add(2)

	// Copy data from source to destination
	go func() {
		defer wg.Done()
		p.copyAndLog(srcConn, dstConn)
	}()

	// Copy data from destination to source
	go func() {
		defer wg.Done()
		p.copyAndLog(dstConn, srcConn)
	}()

	// Wait for both directions to finish copying
	wg.Wait()
}


// resolveForwardAddress resolves the forward address considering custom domain mappings
func (p *TCPProxyService) resolveForwardAddress(forwardAddr string) (*net.TCPAddr, error) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	host, port, err := net.SplitHostPort(forwardAddr)
	if err != nil {
		return nil, err
	}

	// Check if host is a custom domain
	if ip, exists := p.DomainMap[host]; exists {
		host = ip
	}

	return net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
}

// UpdateDomainMapping updates the custom domain mapping in a thread-safe manner
func (p *TCPProxyService) UpdateDomainMapping(domain, ip string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.DomainMap[domain] = ip
}

// copyAndLog copies data between two connections and logs errors if encountered
// func (p *TCPProxyService) copyAndLog(src, dst net.Conn) {
// 	_, err := io.Copy(dst, src)
// 	if err != nil {
// 		// Ignore errors related to closed network connections
// 		if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "use of closed network connection") {
// 			return
// 		}
// 		log.Println("Error copying data:", err)
// 	}
// }
func (p *TCPProxyService) copyAndLog(src, dst net.Conn) {
	_, err := io.Copy(dst, src)

	// Check for specific errors to ignore
	if err != nil {
		// Ignore EOF or "use of closed network connection" errors as they are normal
		if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Connection closed gracefully: %s <-> %s", src.RemoteAddr(), dst.RemoteAddr())
			return
		}

		// Log any other unexpected errors
		log.Printf("Unexpected error copying data from %s to %s: %v", src.RemoteAddr(), dst.RemoteAddr(), err)
	}
}

// GetDomainMappings retrieves all custom domain mappings
func (p *TCPProxyService) GetDomainMappings() map[string]string {
	p.Mu.RLock()
	defer p.Mu.RUnlock()
	return p.DomainMap
}

// DeleteDomainMapping deletes a domain mapping if it exists
func (p *TCPProxyService) DeleteDomainMapping(domain string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if _, exists := p.DomainMap[domain]; !exists {
		return errors.New("domain mapping not found")
	}
	delete(p.DomainMap, domain)
	return nil
}

// GetProxyMappings retrieves all proxy mappings
func (p *TCPProxyService) GetProxyMappings() map[string]models.ProxyMapping {
	p.Mu.RLock()
	defer p.Mu.RUnlock()
	return p.Mappings
}

// UpdateMapping updates the forward address of an existing mapping
func (p *TCPProxyService) UpdateMapping(listenAddr, newForwardAddr string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if _, exists := p.Mappings[listenAddr]; !exists {
		return errors.New("proxy mapping not found")
	}

	// Update the forward address
	p.Mappings[listenAddr] = models.ProxyMapping{ListenAddr: listenAddr, ForwardAddr: newForwardAddr}
	return nil
}

// DeleteMapping deletes a proxy mapping and closes the associated listener
func (p *TCPProxyService) DeleteMapping(listenAddr string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if _, exists := p.Mappings[listenAddr]; !exists {
		return errors.New("proxy mapping not found")
	}

	if listener, ok := p.Listeners[listenAddr]; ok {
		listener.Close()
	}

	delete(p.Mappings, listenAddr)
	delete(p.Listeners, listenAddr)
	return nil
}