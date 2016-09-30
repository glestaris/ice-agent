package ice

import (
	"net/http"
	"net/http/httptest"
)

type handlerFunc func(resp http.ResponseWriter, req *http.Request)

type handler struct {
	method      string
	path        string
	handlerFunc handlerFunc
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

type fakeIceServer struct {
	testServer *httptest.Server
	handlers   []handler
}

func newFakeIceServer() *fakeIceServer {
	return &fakeIceServer{}
}

func (s *fakeIceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *fakeIceServer) handle(method, path string, handlerFunc handlerFunc) {
	s.handlers = append(s.handlers, handler{
		method:      method,
		path:        path,
		handlerFunc: handlerFunc,
	})
}

func (s *fakeIceServer) start() {
	s.testServer = httptest.NewServer(s)
}

func (s *fakeIceServer) endpoint() string {
	return s.testServer.URL
}

func (s *fakeIceServer) stop() {
	s.testServer.Close()
}
