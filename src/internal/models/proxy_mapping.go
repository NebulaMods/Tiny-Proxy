package models

// ProxyMapping represents a mapping from a listening port to a forward address
type ProxyMapping struct {
	ListenAddr  string
	ForwardAddr string
}