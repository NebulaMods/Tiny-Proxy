package models

// ProxyMappingRequest represents a request to add or delete a proxy mapping
type ProxyMappingRequest struct {
	ListenAddr  string `json:"listen_addr"`
	ForwardAddr string `json:"forward_addr"`
}