package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/notjrbauer/fruitvendor"
)

type ProductHandler struct {
	*httprouter.Router

	ProductService fruitvendor.ProductService

	Logger *log.Logger
}

// NewProductHandler returns a new instance of ProductHandler.
func NewProductHandler() *ProductHandler {
	h := &ProductHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.GET("/api/products", h.handleGetProducts)
	h.GET("/api/products/:id", h.handleGetProduct)
	h.POST("/api/products", h.handlePostProduct)
	return h
}

// handleGetProduct handles requests to fetch a single product
func (h *ProductHandler) handleGetProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	p, err := h.ProductService.Product(fruitvendor.ProductID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if p == nil {
		NotFound(w)
	} else {
		encodeJSON(w, &getProductResponse{Product: p}, h.Logger)
	}
}

type getProductResponse struct {
	Product *fruitvendor.Product `json:"product,omitempty"`
	Err     string        `json:"err,omitempty"`
}

// handleGetProducts handles requests to fetch a series of products
func (h *ProductHandler) handleGetProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	p, err := h.ProductService.Products()
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if len(p) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}` + "\n"))
	} else {
		encodeJSON(w, &getProductsResponse{Products: p}, h.Logger)
	}
}

type getProductsResponse struct {
	Products []*fruitvendor.Product `json:"products,omitempty"`
	Err      string          `json:"err,omitempty"`
}

// handlePostProduct handles requests to create a new product
func (h *ProductHandler) handlePostProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request.
	var req postProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	p := req.Product
	p.Token = req.Token
	p.ModTime = time.Time{}

	// Create product.
	switch err := h.ProductService.CreateProduct(p); err {
	case nil:
		encodeJSON(w, &postProductRequest{Product: p}, h.Logger)
	case fruitvendor.ErrProductRequired, fruitvendor.ErrProductIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case fruitvendor.ErrProductExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type postProductRequest struct {
	Product *fruitvendor.Product `json:"product,omitempty"`
	Token   string        `json:"token,omitempty"`
}

type postProductResponse struct {
	Product *fruitvendor.Product `json:"product,omitempty"`
	Err     string        `json:"err,omitempty"`
}

// ProductService represents an HTTP implementation of fruitvendor.ProductService.
type ProductService struct {
	URL *url.URL
}

func (s *ProductService) Product(id fruitvendor.ProductID) (*fruitvendor.Product, error) {
	u := *s.URL
	u.Path = "/api/products/" + url.QueryEscape(string(id))

	// Execute the request.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody getProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Err != "" {
		return nil, fruitvendor.Error(respBody.Err)
	}
	return respBody.Product, nil
}

func (s *ProductService) Products() ([]*fruitvendor.Product, error) {
	u := *s.URL
	u.Path = "/api/products"

	// Execute the request
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody getProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Err != "" {
		return nil, fruitvendor.Error(respBody.Err)
	}
	return respBody.Products, nil
}

func (s *ProductService) CreateProduct(p *fruitvendor.Product) error {
	// Validate arguments.
	if p == nil {
		return fruitvendor.ErrProductRequired
	}

	u := *s.URL
	u.Path = "/api/products"

	// Save token.
	token := p.Token

	reqBody, err := json.Marshal(postProductRequest{Product: p, Token: token})
	if err != nil {
		return err
	}

	// Execute the request.
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody postProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Err != "" {
		return fruitvendor.Error(respBody.Err)
	}

	// Copy returned product.
	*p = *respBody.Product
	p.Token = token

	return err
}

func (s *ProductService) UpdateProduct(id fruitvendor.ProductID, p *fruitvendor.Product) error {
	panic("not implemented")
}

func (s *ProductService) DeleteProduct(id fruitvendor.ProductID, token string) error {
	panic("not implemented")
}
