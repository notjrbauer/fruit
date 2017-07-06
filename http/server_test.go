package http_test

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/notjrbauer/fruit/http"
)

// Now represents the mocked current time.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Server represents a test wrapper for http.Server.
type Server struct {
	*http.Server

	Handler *Handler
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	s := &Server{
		Server:  http.NewServer(),
		Handler: NewHandler(),
	}
	s.Server.Handler = s.Handler.Handler

	// Use random port.
	s.Addr = ":0"

	return s
}

// MustOpenServerClient returns a running server and associated client. Panic on error.
func MustOpenServerClient() (*Server, *http.Client) {
	// Create and open test server.
	s := NewServer()
	if err := s.Open(); err != nil {
		panic(err)
	}

	// Create a client pointing to the server.
	c := http.NewClient()
	c.URL = url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", s.Port())}

	return s, c
}

// VerboseWriter returns a multi-writer to STDERR and w if the "-v" flag is set.
func VerboseWriter(w io.Writer) io.Writer {
	if testing.Verbose() {
		return io.MultiWriter(w, os.Stderr)
	}
	return w
}
