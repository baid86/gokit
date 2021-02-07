// mockHttpServer is Forked from https://git.sr.ht/~ewintr/go-kit/tree/master/test/httpmock.go
// there are minor changes but for most of the code and design remains the same.
// Original License terms https://git.sr.ht/~ewintr/go-kit/tree/master/item/LICENSE

package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
)

// MockResponse represents a response for the mock server to serve
type MockResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	BodyFn     func() []byte
	HeaderFn   func() http.Header
}

type MockServerProcedure struct {
	URI        string
	HTTPMethod string
	Response   MockResponse
	ResponseFn func() MockResponse
}

// MockRecorder provides a way to record request information from every
// successful request.
type MockRecorder interface {
	Record(r *http.Request)
}

// recordedRequest represents recorded structured information about each request
type recordedRequest struct {
	hits     int
	requests []*http.Request
	bodies   [][]byte
}

// MockAssertion represents a common assertion for requests
type MockAssertion struct {
	indexes map[string]int    // indexation for key
	recs    []recordedRequest // request catalog
}

// Record records request hit information
func (m *MockAssertion) Record(r *http.Request) {
	k := m.index(r.URL.Path, r.Method)

	b, _ := ioutil.ReadAll(r.Body)
	if len(b) == 0 {
		b = nil
	}

	if k < 0 {
		m.newIndex(r.URL.Path, r.Method)
		m.recs = append(m.recs, recordedRequest{
			hits:     1,
			requests: []*http.Request{r},
			bodies:   [][]byte{b},
		})
		return
	}

	m.recs[k].hits++
	m.recs[k].requests = append(m.recs[k].requests, r)
	m.recs[k].bodies = append(m.recs[k].bodies, b)
}

// URL  returns the parsed URL so it can be inspected for query params
func (m *MockAssertion) URL(uri, method string) (urlList []*url.URL) {
	k := m.index(uri, method)
	if k < 0 {
		return
	}
	for _, r := range m.recs[k].requests {
		urlList = append(urlList, r.URL)
	}
	return
}

// Hits returns the number of hits for a uri and method
func (m *MockAssertion) Hits(uri, method string) int {
	k := m.index(uri, method)
	if k < 0 {
		return 0
	}

	return m.recs[k].hits
}

// Headers returns a slice of request headers
func (m *MockAssertion) Headers(uri, method string) []http.Header {
	k := m.index(uri, method)
	if k < 0 {
		return nil
	}

	headers := make([]http.Header, len(m.recs[k].requests))
	for i, r := range m.recs[k].requests {

		// remove default headers
		if _, ok := r.Header["Content-Length"]; ok {
			r.Header.Del("Content-Length")
		}

		if v, ok := r.Header["User-Agent"]; ok {
			if reflect.DeepEqual([]string{"Go-http-client/1.1"}, v) {
				r.Header.Del("User-Agent")
			}
		}

		if v, ok := r.Header["Accept-Encoding"]; ok {
			if reflect.DeepEqual([]string{"gzip"}, v) {
				r.Header.Del("Accept-Encoding")
			}
		}

		if len(r.Header) == 0 {
			continue
		}

		headers[i] = r.Header
	}
	return headers
}

// Body returns request body
func (m *MockAssertion) Body(uri, method string) [][]byte {
	k := m.index(uri, method)
	if k < 0 {
		return nil
	}

	return m.recs[k].bodies
}

// Reset sets all unexpected properties to their zero value
func (m *MockAssertion) Reset() error {
	m.indexes = make(map[string]int)
	m.recs = make([]recordedRequest, 0)
	return nil
}

// index indexes a key composed of the uri and method and returns the position
// for this key in a list if it was indexed before.
func (m *MockAssertion) index(uri, method string) int {
	if len(m.indexes) == 0 {
		m.indexes = make(map[string]int)
	}

	k := strings.ToLower(uri + method)

	if i, ok := m.indexes[k]; ok {
		return i
	}

	return -1
}

func (m *MockAssertion) newIndex(uri, method string) int {
	k := strings.ToLower(uri + method)
	m.indexes[k] = len(m.indexes)
	return m.indexes[k]
}

// NewMockServer return a mock HTTP server to test requests
func NewMockServer(rec MockRecorder, procedures ...MockServerProcedure) *httptest.Server {
	var handler http.Handler

	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			for _, proc := range procedures {
				if proc.URI == r.URL.Path && proc.HTTPMethod == r.Method {

					headers := w.Header()
					procResponse := proc.Response
					if proc.ResponseFn != nil {
						procResponse = proc.ResponseFn()
					}
					for hkey, hvalue := range procResponse.Headers {
						headers[hkey] = hvalue
					}

					code := procResponse.StatusCode
					if code == 0 {
						code = http.StatusOK
					}

					w.WriteHeader(code)
					if procResponse.BodyFn != nil {
						w.Write(procResponse.BodyFn())
					} else {
						w.Write(procResponse.Body)
					}

					if rec != nil {
						rec.Record(r)
					}
					return
				}
			}

			w.WriteHeader(http.StatusNotFound)
			return
		})

	return httptest.NewServer(handler)
}
