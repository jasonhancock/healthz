package healthz

import (
	"context"
	"fmt"
	"testing"

	"github.com/cheekybits/is"
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
	is := is.New(t)

	count = 0

	c := testCheck{}
	ctx := context.Background()

	r := c.Check(ctx)
	is.Equal(count, 1)
	is.Equal(r.Metadata["count"], "1")

	r = c.Check(ctx)
	is.Equal(count, 2)
	is.Equal(r.Metadata["count"], "2")
}
