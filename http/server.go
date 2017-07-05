package http

import (
	"net"
	"net/http"
	"net/url"

	"github.com/notjrbauer/fruitvendor"
)

// DefaultAddr is the default bind address
const DefaultAddr = ":3000"

// Server represents a HTTP server.
type Server struct {
	ln net.Listener

	// Handler to server.
	Handler *Handler

	// Bind address to open.
	Addr string
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	return &Server{
		Addr: DefaultAddr,
	}
}

// Open opens a socket and servers the HTTP server.
func (s *Server) Open() error {
	// Open socket.

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.ln = ln

	// Start HTTP server.
	go func() { http.Serve(s.ln, s.Handler) }()
	//http.Serve(s.ln, s.Handler)

	return nil
}

// Close closes the socket
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port returns the port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

// Client represents a client to connect to the HTTP server.
type Client struct {
	URL            url.URL
	productService ProductService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{}
	c.productService.URL = &c.URL
	return c
}

func (c *Client) ProductService() fruit.ProductService {
	return &c.productService
}
