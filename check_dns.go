package healthz

import (
	"context"
	"net"
)

// CheckDNS attempts to resolve a given domain
type CheckDNS struct {
	address string
	network string
}

// NewCheckDNS creates a new DNS check.
func NewCheckDNS(network, address string) CheckDNS {
	return CheckDNS{
		address: address,
		network: network,
	}
}

// Check is called by the checker and attempts to resolve the address
func (c CheckDNS) Check(ctx context.Context) *Response {
	ip, err := net.ResolveIPAddr(c.network, c.address)

	var addr string
	if ip != nil {
		addr = ip.String()
	}

	return &Response{
		Error: err,
		Metadata: map[string]string{
			"ip": addr,
		},
	}
}
