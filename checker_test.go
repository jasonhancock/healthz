package healthz

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type sleep struct {
	duration time.Duration
}

func newSleep(dur time.Duration) sleep {
	return sleep{
		duration: dur,
	}
}

func (c sleep) Check(ctx context.Context) *Response {
	time.Sleep(c.duration)
	return &Response{}
}

func TestChecker(t *testing.T) {
	c := NewChecker(WithTimeout(100 * time.Millisecond))
	err := c.AddCheck("test", testCheck{})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "http://127.0.0.1/healthz", nil)
	w := httptest.NewRecorder()

	c.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	m, err := parseResponse(w.Body)
	require.NoError(t, err)

	v, ok := m["test"]
	require.True(t, ok)
	require.Len(t, v.Metadata, 1)
	require.Equal(t, "", v.ErrorMessage)

	// Add another check, one designed to fail
	require.NoError(t, c.AddCheck("sleep", newSleep(200*time.Millisecond)))

	w = httptest.NewRecorder()

	req = httptest.NewRequest("GET", "http://127.0.0.1/healthz", nil)
	w = httptest.NewRecorder()

	c.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)

	m, err = parseResponse(w.Body)
	require.NoError(t, err)

	v, ok = m["sleep"]
	require.True(t, ok)
	require.Equal(t, ErrorCheckTimedOut.Error(), v.ErrorMessage)
}

func TestCheckerEmpty(t *testing.T) {
	c := NewChecker(WithTimeout(100 * time.Millisecond))

	req := httptest.NewRequest("GET", "http://127.0.0.1/healthz", nil)
	w := httptest.NewRecorder()

	c.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCheckerUnknownCheck(t *testing.T) {
	c := NewChecker(WithTimeout(100 * time.Millisecond))
	require.NoError(t, c.AddCheck("test", testCheck{}))

	req := httptest.NewRequest("GET", "http://127.0.0.1/healthz/foo", nil)
	w := httptest.NewRecorder()

	c.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func parseResponse(r io.Reader) (map[string]Response, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	m := make(map[string]Response)
	if err = json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	return m, nil
}
