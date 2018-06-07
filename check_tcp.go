package healthz

import (
	"context"
	"net"
	"time"
)

// CheckTCP is a TCP health check
type CheckTCP struct {
	addr    string
	timeout time.Duration
}

// NewCheckTCP creates a new TCP check.
func NewCheckTCP(addr string, timeout time.Duration) CheckTCP {
	return CheckTCP{
		addr:    addr,
		timeout: timeout,
	}
}

// Check is called by the checker and attempts to connect to addr.
func (c CheckTCP) Check(ctx context.Context) *Response {
	d := &net.Dialer{Timeout: c.timeout}
	conn, err := d.Dial("tcp", c.addr)
	if err == nil {
		conn.Close()
	}
	return &Response{
		Error: err,
	}
}
