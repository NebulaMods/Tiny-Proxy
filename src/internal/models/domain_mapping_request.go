package models

// DomainMappingRequest represents a request to update or delete a domain mapping
type DomainMappingRequest struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
}