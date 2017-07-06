package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/notjrbauer/fruit"
)

const ErrInvalidJSON = fruit.Error("invalid json")

// Handler is a collection of all the service handlers.

type Handler struct {
	ProductHandler *ProductHandler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/products") {
		h.ProductHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = fruit.ErrInternal
	}

	// Write generic response.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

// errorResponse is a generic response for sending an error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// encodeJson encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

// NotFound writes an API error message to the response.
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}`))
}
