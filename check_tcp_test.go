package healthz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCheckTCP(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()

	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	c := NewCheckTCP(u.Host, 5*time.Second)
	result := c.Check(context.Background())
	require.NoError(t, result.Error)

	// Shut the server down, expect an error
	server.Close()
	result = c.Check(context.Background())
	require.Error(t, result.Error)
}

func ExampleCheckTCP() {
	checker := NewChecker()
	checker.AddCheck("app", NewCheckTCP("127.0.0.1:22", 5*time.Second))
	http.ListenAndServe(":8080", checker)
}
