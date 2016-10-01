package testing

import (
	"net/http"
	"net/http/httptest"
)

// HandlerFunc is the signature of any function used to handle HTTP requests in
// the test server.
type HandlerFunc func(resp http.ResponseWriter, req *http.Request)

type handler struct {
	method      string
	path        string
	handlerFunc HandlerFunc
}

func (h handler) matches(req *http.Request) bool {
	if req.Method != h.method {
		return false
	}
	if req.URL.Path != h.path {
		return false
	}
	return true
}

func (h handler) handle(resp http.ResponseWriter, req *http.Request) {
	h.handlerFunc(resp, req)
}

// FakeIceServer is a test server that replicates an iCE registry server
// interface.
type FakeIceServer struct {
	testServer *httptest.Server
	handlers   []handler
}

// NewFakeIceServer is a constructor that provides a reference to a
// FakeIceServer instance.
func NewFakeIceServer() *FakeIceServer {
	return &FakeIceServer{}
}

func (s *FakeIceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	for _, h := range s.handlers {
		if h.matches(req) {
			h.handle(resp, req)
			return
		}
	}

	resp.WriteHeader(http.StatusInternalServerError)
	_, err := resp.Write([]byte("Could not find handler for this request\n"))
	if err != nil {
		panic(err)
	}
}

// Handle registers a HandleFunc for requests directed to the provided path and
// HTTP verb.
func (s *FakeIceServer) Handle(method, path string, handlerFunc HandlerFunc) {
	s.handlers = append(s.handlers, handler{
		method:      method,
		path:        path,
		handlerFunc: handlerFunc,
	})
}

// Start starts the test server.
func (s *FakeIceServer) Start() {
	s.testServer = httptest.NewServer(s)
}

// Endpoint returns the URL of the test server. The test server picks a random
// port every time it starts.
func (s *FakeIceServer) Endpoint() string {
	return s.testServer.URL
}

// Stop kills the connections to the test server and stops the server.
func (s *FakeIceServer) Stop() {
	s.testServer.Close()
}
