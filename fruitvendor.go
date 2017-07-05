package fruitvendor

import "time"

type ProductID string

type Product struct {
	ID          ProductID `json:"productID" storm:"id"`
	Token       string    `json:"-"`
	Name        string    `json:"name, omitempty"`
	SKU         string    `json:"sku"`
	Type        string    `json:"type"`
	Color       string    `json:"color"`
	Description string    `json:"description, omitempty"`
	ModTime     time.Time `json:"modTime"`
}

// Client creates a connection to the services.
type Client interface {
	ProductService() ProductService
}

// ProductService represents a service for managing products
type ProductService interface {
	Product(id ProductID) (*Product, error)
	Products() ([]*Product, error)
	CreateProduct(p *Product) error
	UpdateProduct(id ProductID, p *Product) error
	DeleteProduct(id ProductID, token string) error
}
