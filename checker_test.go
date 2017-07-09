package healthz

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cheekybits/is"
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
	is := is.New(t)

	c := NewChecker(WithTimeout(100 * time.Millisecond))
	err := c.AddCheck("test", testCheck{})
	is.NoErr(err)

	req := httptest.NewRequest("GET", "http://127.0.0.1/healthz", nil)
	w := httptest.NewRecorder()

	c.ServeHTTP(w, req)
	is.Equal(w.Code, 200)

	m, err := parseResponse(w.Body)
	is.NoErr(err)

	v, ok := m["test"]
	is.OK(ok)
	is.Equal(1, len(v.Metadata))
	is.Equal("", v.ErrorMessage)

	// Add another check, one designed to fail
	err = c.AddCheck("sleep", newSleep(200*time.Millisecond))
	is.NoErr(err)

	w = httptest.NewRecorder()

	req = httptest.NewRequest("GET", "http://127.0.0.1/healthz", nil)
	w = httptest.NewRecorder()

	c.ServeHTTP(w, req)
	is.Equal(w.Code, 500)

	m, err = parseResponse(w.Body)
	is.NoErr(err)

	v, ok = m["sleep"]
	is.OK(ok)
	is.Equal(ErrorCheckTimedOut.Error(), v.ErrorMessage)
}

func TestCheckerEmpty(t *testing.T) {
	is := is.New(t)

	c := NewChecker(WithTimeout(100 * time.Millisecond))

	req := httptest.NewRequest("GET", "http://127.0.0.1/healthz", nil)
	w := httptest.NewRecorder()

	c.ServeHTTP(w, req)
	is.Equal(w.Code, 200)
}

func TestCheckerUnknownCheck(t *testing.T) {
	is := is.New(t)

	c := NewChecker(WithTimeout(100 * time.Millisecond))
	err := c.AddCheck("test", testCheck{})
	is.NoErr(err)

	req := httptest.NewRequest("GET", "http://127.0.0.1/healthz/foo", nil)
	w := httptest.NewRecorder()

	c.ServeHTTP(w, req)
	is.Equal(w.Code, 404)
}

func parseResponse(r io.Reader) (map[string]Response, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	m := make(map[string]Response)
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
