package bolt

import (
	"time"

	"github.com/notjrbauer/fruitvendor"
)

type ProductService struct {
	client *Client
}

// Product returns a product by ID.
func (s *ProductService) Product(id fruitvendor.ProductID) (*fruitvendor.Product, error) {
	// Start read-only transaction.
	tx, err := s.client.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal product.
	var p fruitvendor.Product
	products := s.client.db.From("Products")

	if err := products.One("ID", id, &p); err != nil {
		return nil, err
	} else if &p == nil {
		return nil, nil
	}

	return &p, nil
}

func (s *ProductService) Products() ([]*fruitvendor.Product, error) {
	var products []*fruitvendor.Product
	if err := s.client.db.From("Products").All(&products); err != nil {
		return nil, err
	}
	return products, nil
}

// CreateProduct creates a new product.
func (s *ProductService) CreateProduct(p *fruitvendor.Product) error {
	// Require id
	if p.ID == "" {
		return fruitvendor.ErrProductIDRequired
	}

	bucket := s.client.db.From("Products")
	// Start the read-write transaction.
	tx, err := bucket.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify product doesn't already exist.

	var product fruitvendor.Product
	tx.One("ID", p.ID, &product)

	if product.ID != "" {
		return fruitvendor.ErrProductExists
	}

	// Update modified time.
	p.ModTime = time.Now().UTC()

	if err := tx.Save(p); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateProduct updates an existing product.
func (s *ProductService) UpdateProduct(id fruitvendor.ProductID, p *fruitvendor.Product) error {
	// Start read-write transaction.

	tx, err := s.client.db.From("Products").Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Find record.
	var product fruitvendor.Product
	if err := tx.One("ID", p.ID, &product); err != nil {
		return fruitvendor.ErrProductNotFound
	}

	// Apply changes.
	var d fruitvendor.Product
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
func (s *ProductService) DeleteProduct(id fruitvendor.ProductID, token string) error {
	// Start the read-write transaction.
	tx, err := s.client.db.From("Products").Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Find record.
	var product fruitvendor.Product
	if err := tx.One("ID", id, &product); err != nil {
		return fruitvendor.ErrProductNotFound
	}

	if err := tx.DeleteStruct(&product); err != nil {
		return err
	}

	return tx.Commit()
}
