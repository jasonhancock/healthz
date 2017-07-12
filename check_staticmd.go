package healthz

import (
	"context"
)

// CheckStaticMD is a static metadata healthz check
type CheckStaticMD struct {
	md map[string]string
}

// NewCheckStaticMD creates a new static metadata check. The metadata passed in here will be
// returned on every check
func NewCheckStaticMD(md map[string]string) CheckStaticMD {
	return CheckStaticMD{
		md: md,
	}
}

// Check is called by the checker and returns the metadata
func (c CheckStaticMD) Check(ctx context.Context) *Response {
	return &Response{
		Metadata: c.md,
	}
}
