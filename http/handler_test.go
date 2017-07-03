package http_test

import "github.com/notjrbauer/caps/http"

// Handler represents a test wrapper for http.Handler.
type Handler struct {
	*http.Handler

	ProductHandler *ProductHandler
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler:        &http.Handler{},
		ProductHandler: NewProductHandler(),
	}
	h.Handler.ProductHandler = h.ProductHandler.ProductHandler
	return h
}
