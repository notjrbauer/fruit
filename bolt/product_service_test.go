package bolt_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/notjrbauer/fruitvendor"
)

func TestProductService_CreateProduct(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ProductService()

	product := fruitvendor.Product{
		ID:          "ID",
		SKU:         "SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	if err := s.CreateProduct(&product); err != nil {
		t.Fatal(err)
	}

	other, err := s.Product("ID")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&product, other) {
		t.Fatalf("unexpected product: %+v", other)
	}
}

func TestProductService_CreateProduct_ErrProductIDRequired(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	product := fruitvendor.Product{
		ID:          "",
		SKU:         "SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	if err := c.ProductService().CreateProduct(&product); err != fruitvendor.ErrProductIDRequired {
		t.Fatalf("expected error with without id: %+v", product)
	}
}

func TestProductService_CreateProduct_ErrProductExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	if err := c.ProductService().CreateProduct(&fruitvendor.Product{ID: "X"}); err != nil {
		t.Fatal(err)
	}

	if err := c.ProductService().CreateProduct(&fruitvendor.Product{ID: "X"}); err != fruitvendor.ErrProductExists {
		t.Fatal(errors.New("expected error when creating same product"))
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ProductService()

	// Create new product.
	product := fruitvendor.Product{
		ID:          "XXX",
		SKU:         "OLD_SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	if err := s.CreateProduct(&product); err != nil {
		t.Fatal(err)
	}

	product.SKU = "NEW_SKU"

	// Update product
	if err := s.UpdateProduct("XXX", &product); err != nil {
		t.Fatal(err)
	}

	// Verify product updated.
	if p, err := s.Product(product.ID); err != nil {
		t.Fatal(err)
	} else if p.SKU != "NEW_SKU" {
		t.Fatalf("unexpected product sku: %s", p.SKU)
	}
}

func TestProductService_UpdateProduct_ErrProductNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ProductService()

	// Create new product.
	product := fruitvendor.Product{
		ID:          "XXX",
		SKU:         "OLD_SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	// Update product
	if err := s.UpdateProduct("XXX", &product); err != fruitvendor.ErrProductNotFound {
		t.Fatal(err)
	}
}

func TestProductService_Update_ErrProductDoesNotExist(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	if err := c.ProductService().UpdateProduct("XXX", &fruitvendor.Product{ID: "X"}); err == nil {
		t.Fatal("product should not update non-existing product")
	}
}

func TestProductService_Delete(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ProductService()

	// Create new product.
	product := fruitvendor.Product{
		ID:          "XXX",
		SKU:         "OLD_SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	if err := s.CreateProduct(&product); err != nil {
		t.Fatal(err)
	}

	// Delete product.
	if err := s.DeleteProduct(product.ID, product.Token); err != nil {
		t.Fatal(err)
	}

	// Verify removal of product..
	if _, err := s.Product(product.ID); err == nil {
		t.Fatal(errors.New("product was not removed"))
	}
}

func TestProductService_Delete_ErrProductNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ProductService()

	// Create new product.
	product := fruitvendor.Product{
		ID:          "XXX",
		SKU:         "OLD_SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	// Delete product.
	if err := s.DeleteProduct(product.ID, product.Token); err != fruitvendor.ErrProductNotFound {
		t.Fatalf("expected error with non-existing product: %+v", product)
	}
}

func TestProductService_Products(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ProductService()

	// Create new product.
	product := fruitvendor.Product{
		ID:          "XXX",
		SKU:         "OLD_SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	if err := s.CreateProduct(&product); err != nil {
		t.Fatal(err)
	}

	// Create second product.
	product = fruitvendor.Product{
		ID:          "YYY",
		SKU:         "OLD_SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	if err := s.CreateProduct(&product); err != nil {
		t.Fatal(err)
	}

	// Fetch products.
	products, err := s.Products()
	if err != nil {
		t.Fatal(err)
	}

	// Verify removal of product..
	if len(products) != 2 {
		t.Fatal(errors.New("multiple products were not returned"))
	}
}

func TestProductService_Products_Empty(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	if p, _ := c.ProductService().Products(); p == nil {
		t.Fatal("expected empty product array")
	}
}
