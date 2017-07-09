package healthz

import (
	"context"
	//	"encoding/json"
	//	"log"
	//	"net/http"
)

// Check is the interface for a healthz check. A Check has a Check function that
// takes a context and returns a Response
type Check interface {
	Check(context.Context) *Response
}

// Response is the response from a Check
type Response struct {
	Metadata     map[string]string `json:"metadata,omitempty"`
	Error        error             `json:"-"`
	ErrorMessage string            `json:"error,omitempty"`
}

// NewResponse returns a new *Response with an initialized Metadata map
func NewResponse() *Response {
	return &Response{
		Metadata: make(map[string]string),
	}
}
