package healthz

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckDNS(t *testing.T) {
	var tests = []struct {
		address     string
		network     string
		expectError bool
	}{
		{"www.google.com", "ip4", false},
		{"www.google.com", "ip6", false},
		{"error.jasonhancock.com", "ip4", true},
		{"jasonhancock.com", "ip4", false},
		{"jasonhancock.com", "ip6", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s", tt.address, tt.network), func(t *testing.T) {
			c := NewCheckDNS(tt.network, tt.address)
			result := c.Check(context.Background())
			if tt.expectError {
				require.Error(t, result.Error)
			} else {
				require.NoError(t, result.Error)
			}
		})
	}
}

func ExampleCheckDNS() {
	checker := NewChecker()
	checker.AddCheck("dns-www.google.com", NewCheckDNS("ip4", "www.google.com"))
	http.ListenAndServe(":8080", checker)
}
