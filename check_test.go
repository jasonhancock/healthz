package healthz

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var count int

type testCheck struct{}

func (c testCheck) Check(ctx context.Context) *Response {
	count++
	return &Response{
		Metadata: map[string]string{
			"count": fmt.Sprintf("%d", count),
		},
	}
}

func TestCheck(t *testing.T) {
	count = 0

	c := testCheck{}
	ctx := context.Background()

	r := c.Check(ctx)
	require.Equal(t, 1, count)
	require.Equal(t, "1", r.Metadata["count"])

	r = c.Check(ctx)
	require.Equal(t, 2, count)
	require.Equal(t, "2", r.Metadata["count"])
}
