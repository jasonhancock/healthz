package healthz

import (
	"context"
)

// StaticMD is a static metadata healthz check
type StaticMD struct {
	md map[string]string
}

// NewStaticMD creates a new static metadata check. The metadata passed in here will be
// returned on every check
func NewStaticMD(md map[string]string) StaticMD {
	return StaticMD{
		md: md,
	}
}

// Check is called by the checker and returns the metadata
func (c StaticMD) Check(ctx context.Context) *Response {
	return &Response{
		Metadata: c.md,
	}
}
