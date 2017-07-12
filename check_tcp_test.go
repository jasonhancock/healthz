package healthz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cheekybits/is"
)

func TestCheckTCP(t *testing.T) {
	is := is.New(t)

	server := httptest.NewServer(nil)
	defer server.Close()

	u, err := url.Parse(server.URL)
	is.NoErr(err)

	c := NewCheckTCP(u.Host, 5*time.Second)
	result := c.Check(context.Background())
	is.NoErr(result.Error)

	// Shut the server down, expect an error
	server.Close()
	result = c.Check(context.Background())
	is.Err(result.Error)
}

func ExampleCheckTCP() {
	checker := NewChecker()
	checker.AddCheck("app", NewCheckTCP("127.0.0.1:22", 5*time.Second))
	http.ListenAndServe(":8080", checker)
}
