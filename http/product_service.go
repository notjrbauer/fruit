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

	ProductService fruit.ProductService

	Logger *log.Logger
}

// NewProductHandler returns a new instance of ProductHandler.
func NewProductHandler() *ProductHandler {
	h := &ProductHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.GET("/api/products", h.handleGetProducts)
	h.POST("/api/products", h.handlePostProduct)
	h.PUT("/api/products", h.handlePutProduct)
	h.DELETE("/api/products", h.handleDeleteProduct)

	h.GET("/api/products/:id", h.handleGetProduct)
	return h
}

// handleGetProduct handles requests to fetch a single product
func (h *ProductHandler) handleGetProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	p, err := h.ProductService.Product(fruit.ProductID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if p == nil {
		NotFound(w)
	} else {
		encodeJSON(w, &getProductResponse{Product: p}, h.Logger)
	}
}

type getProductResponse struct {
	Product *fruit.Product `json:"product,omitempty"`
	Err     string         `json:"err,omitempty"`
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
	Products []*fruit.Product `json:"products,omitempty"`
	Err      string           `json:"err,omitempty"`
}

// handlePostProduct handles requests to create a new product.
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
	case fruit.ErrProductRequired, fruit.ErrProductIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case fruit.ErrProductExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type postProductRequest struct {
	Product *fruit.Product `json:"product,omitempty"`
	Token   string         `json:"token,omitempty"`
}

type postProductResponse struct {
	Product *fruit.Product `json:"product,omitempty"`
	Err     string         `json:"err,omitempty"`
}

// handlePutProduct handles requests to update a product.
func (h *ProductHandler) handlePutProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request.
	var req putProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	p := req.Product
	p.ID = req.ID
	p.ModTime = time.Time{}

	// Create product.
	// TODO: Add Token
	switch err := h.ProductService.UpdateProduct(p.ID, p); err {
	case nil:
		encodeJSON(w, &putProductResponse{Product: p}, h.Logger)
	case fruit.ErrProductRequired, fruit.ErrProductIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case fruit.ErrProductNotFound:
		Error(w, err, http.StatusNotFound, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type putProductRequest struct {
	Product *fruit.Product  `json:"product,omitempty"`
	ID      fruit.ProductID `json:"id,omitempty"`
}

type putProductResponse struct {
	Product *fruit.Product `json:"product,omitempty"`
	Err     string         `json:"err,omitempty"`
}

// handleDeleteProduct handles requests to update a product.
func (h *ProductHandler) handleDeleteProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request.
	var req deleteProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	// Delete product.
	switch err := h.ProductService.DeleteProduct(req.ID, req.Token); err {
	case nil:
		encodeJSON(w, &deleteProductResponse{}, h.Logger)
	case fruit.ErrProductNotFound:
		Error(w, err, http.StatusNotFound, h.Logger)
	case fruit.ErrProductRequired, fruit.ErrProductIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type deleteProductRequest struct {
	ID    fruit.ProductID `json:"id,omitempty"`
	Token string          `json:"token"`
}

type deleteProductResponse struct {
	Err string `json:"err,omitempty"`
}

// ProductService represents an HTTP implementation of fruit.ProductService.
type ProductService struct {
	URL *url.URL
}

func (s *ProductService) Product(id fruit.ProductID) (*fruit.Product, error) {
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
		return nil, fruit.Error(respBody.Err)
	}
	return respBody.Product, nil
}

func (s *ProductService) Products() ([]*fruit.Product, error) {
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
		return nil, fruit.Error(respBody.Err)
	}
	return respBody.Products, nil
}

func (s *ProductService) CreateProduct(p *fruit.Product) error {
	// Validate arguments.
	if p == nil {
		return fruit.ErrProductRequired
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
		return fruit.Error(respBody.Err)
	}

	// Copy returned product.
	*p = *respBody.Product
	p.Token = token

	return err
}

func (s *ProductService) UpdateProduct(id fruit.ProductID, p *fruit.Product) error {
	// Validate arguments.
	if id == "" {
		return fruit.ErrProductIDRequired
	}

	u := *s.URL
	u.Path = "/api/products"

	reqBody, err := json.Marshal(putProductRequest{Product: p, ID: id})
	if err != nil {
		return err
	}

	// Create request.
	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// Execute request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody putProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Err != "" {
		return fruit.Error(respBody.Err)
	}

	// Copy returned product.
	// TODO: Remove ability to generate ID
	*p = *respBody.Product
	p.ID = id
	return nil
}

func (s *ProductService) DeleteProduct(id fruit.ProductID, token string) error {
	// Validate arguments.
	if id == "" {
		return fruit.ErrProductIDRequired
	}

	u := *s.URL
	u.Path = "/api/products"

	reqBody, err := json.Marshal(deleteProductRequest{ID: id, Token: token})
	if err != nil {
		return err
	}

	// Create request.
	req, err := http.NewRequest(http.MethodDelete, u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// Execute request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody deleteProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Err != "" {
		return fruit.Error(respBody.Err)
	}

	return nil
}
