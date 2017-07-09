package hzchecks

import (
	"context"
	"testing"

	"github.com/cheekybits/is"
)

func TestStaticMD(t *testing.T) {
	is := is.New(t)

	md := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	smd := NewStaticMD(md)
	result := smd.Check(context.Background())

	is.Equal(2, len(result.Metadata))
	value1, ok := result.Metadata["key1"]
	is.OK(ok)
	is.Equal("value1", value1)
}

// TODO: include an example