package healthz

import (
	"context"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckStaticMD(t *testing.T) {
	md := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	c := NewCheckStaticMD(md)
	result := c.Check(context.Background())

	require.Len(t, result.Metadata, 2)
	value1, ok := result.Metadata["key1"]
	require.True(t, ok)
	require.Equal(t, "value1", value1)
}

func ExampleCheckStaticMD() {
	checker := NewChecker()
	checker.AddCheck("app", NewCheckStaticMD(map[string]string{
		"go_version": runtime.Version(),
		"go_arch":    runtime.GOARCH,
	}))
	http.ListenAndServe(":8080", checker)
}
