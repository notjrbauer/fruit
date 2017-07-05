package bolt

import (
	"time"

	"github.com/notjrbauer/fruitvendor"
)

type ProductService struct {
	client *Client
}

// Product returns a product by ID.
func (s *ProductService) Product(id fruit.ProductID) (*fruit.Product, error) {
	// Start read-only transaction.
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal product.
	var p fruit.Product
	products := s.client.db.From("Products")

	if err := products.One("ID", id, &p); err != nil {
		return nil, err
	} else if &p == nil {
		return nil, nil
	}

	return &p, nil
}

func (s *ProductService) Products() ([]*fruit.Product, error) {
	var products []*fruit.Product
	if err := s.client.db.From("Products").All(&products); err != nil {
		return nil, err
	}
	return products, nil
}

// CreateProduct creates a new product.
func (s *ProductService) CreateProduct(p *fruit.Product) error {
	// Require id
	if p.ID == "" {
		return fruit.ErrProductIDRequired
	}

	bucket := s.client.db.From("Products")
	// Start the read-write transaction.
	tx, err := bucket.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify product doesn't already exist.

	var product fruit.Product
	tx.One("ID", p.ID, &product)

	if product.ID != "" {
		return fruit.ErrProductExists
	}

	// Update modified time.
	p.ModTime = time.Now().UTC()

	if err := tx.Save(p); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateProduct updates an existing product.
func (s *ProductService) UpdateProduct(id fruit.ProductID, p *fruit.Product) error {
	// Start read-write transaction.

	tx, err := s.client.db.From("Products").Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Find record.
	var product fruit.Product
	if err := tx.One("ID", p.ID, &product); err != nil {
		return fruit.ErrProductNotFound
	}

	// Apply changes.
	var d fruit.Product
	d.ID = p.ID
	d.Color = p.Color
	d.Description = p.Description
	d.Name = p.Name
	d.SKU = p.SKU
	d.Type = p.Type
	d.ModTime = time.Now().UTC()

	if err := tx.Update(&d); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteProduct removes an existing product.
func (s *ProductService) DeleteProduct(id fruit.ProductID, token string) error {
	// Start the read-write transaction.
	tx, err := s.client.db.From("Products").Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Find record.
	var product fruit.Product
	if err := tx.One("ID", id, &product); err != nil {
		return fruit.ErrProductNotFound
	}

	if err := tx.DeleteStruct(&product); err != nil {
		return err
	}

	return tx.Commit()
}
