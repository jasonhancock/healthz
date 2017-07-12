package healthz

import (
	"context"
	"net/http"
	"runtime"
	"testing"

	"github.com/cheekybits/is"
)

func TestCheckStaticMD(t *testing.T) {
	is := is.New(t)

	md := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	c := NewCheckStaticMD(md)
	result := c.Check(context.Background())

	is.Equal(2, len(result.Metadata))
	value1, ok := result.Metadata["key1"]
	is.OK(ok)
	is.Equal("value1", value1)
}

func ExampleCheckStaticMD() {
	checker := NewChecker()
	checker.AddCheck("app", NewCheckStaticMD(map[string]string{
		"go_version": runtime.Version(),
		"go_arch":    runtime.GOARCH,
	}))
	http.ListenAndServe(":8080", checker)
}
