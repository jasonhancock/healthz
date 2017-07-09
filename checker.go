package healthz

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ErrorCheckTimedOut is the error returned when an check times out
var ErrorCheckTimedOut = errors.New("healthz check timed out")

type options struct {
	timeout time.Duration
	prefix  string
}

// Checker is the responsible for running all of the checks
type Checker struct {
	checks  map[string]Check
	timeout time.Duration
	lock    sync.RWMutex
	prefix  string
}

// CheckerOption is used to customize the checker
type CheckerOption func(*options)

// NewChecker creates a new Checker. It can take a set of options to customize the timeout or URL prefix
func NewChecker(opts ...CheckerOption) *Checker {
	opt := &options{
		timeout: 10 * time.Second,
		prefix:  "/healthz",
	}
	for _, o := range opts {
		o(opt)
	}

	return &Checker{
		timeout: opt.timeout,
		prefix:  opt.prefix,
		checks:  make(map[string]Check),
	}
}

// WithTimeout sets the checker timeout. Default is 10s.
func WithTimeout(timeout time.Duration) CheckerOption {
	return func(c *options) {
		c.timeout = timeout
	}
}

// WithPrefix sets the path prefix. Default is /healthz
func WithPrefix(prefix string) CheckerOption {
	return func(c *options) {
		c.prefix = prefix
	}
}

// AddCheck adds the specified check to the list of checks that will be run by the checker when hitting the
// prefix endpoint. The individual check can be called by hitting the "<prefix>/<name>" endpoint.
func (c *Checker) AddCheck(name string, chk Check) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.checks[name]; ok {
		return errors.New("healthz check with name already exists: " + name)
	}

	c.checks[name] = chk
	return nil
}

// Prefix returns the configured url prefix for the Checker
func (c *Checker) Prefix() string {
	return c.prefix
}

type namedResponse struct {
	name   string
	result *Response
}

func (c *Checker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	ctx, cancel := context.WithTimeout(r.Context(), c.timeout)
	defer cancel()

	var checks map[string]Check

	if r.URL.Path != c.prefix {
		name := strings.TrimPrefix(r.URL.Path, c.prefix+"/")
		c.lock.RLock()
		if _, ok := c.checks[name]; !ok {
			c.lock.RUnlock()
			http.Error(w, "No such healthz check", http.StatusNotFound)
			return
		}
		c.lock.RUnlock()
		checks = make(map[string]Check, 1)
		checks[name] = c.checks[name]
	} else {
		c.lock.RLock()
		checks = make(map[string]Check, len(c.checks))
		for k, v := range c.checks {
			checks[k] = v
		}
		c.lock.RUnlock()
	}

	numChecks := len(checks)
	resp := make(map[string]*Response, numChecks)
	timedout := make(map[string]struct{}, numChecks)

	res := make(chan *namedResponse)
	for name, check := range checks {
		timedout[name] = struct{}{}
		go func(chk Check, name string) {
			r := &namedResponse{
				name:   name,
				result: chk.Check(ctx),
			}

			select {
			case <-ctx.Done():
				return
			case res <- r:
			}
		}(check, name)
	}

	for i := 0; i < numChecks; i++ {
		select {
		case <-ctx.Done():
			for name := range timedout {
				resp[name] = &Response{
					Error:        ErrorCheckTimedOut,
					ErrorMessage: ErrorCheckTimedOut.Error(),
				}
				status = http.StatusInternalServerError
			}
			break
		case r := <-res:
			if r.result.Error != nil {
				r.result.ErrorMessage = r.result.Error.Error()
				status = http.StatusInternalServerError
			}
			resp[r.name] = r.result
			delete(timedout, r.name)
		}
	}

	if r.Method != http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		data, err := json.MarshalIndent(&resp, "", "\t")
		if err != nil {
			log.Println(err)
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)
		w.Write(data)
	} else {
		w.WriteHeader(status)
	}
}
