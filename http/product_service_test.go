package http_test

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/notjrbauer/fruitvendor"
	"github.com/notjrbauer/fruitvendor/http"
	"github.com/notjrbauer/fruitvendor/mock"
)

// ProductHandler represents a test wrapper for http.ProductHandler
type ProductHandler struct {
	*http.ProductHandler

	ProductService mock.ProductService
	LogOutput      bytes.Buffer
}

func NewProductHandler() *ProductHandler {
	h := &ProductHandler{ProductHandler: http.NewProductHandler()}
	h.ProductHandler.ProductService = &h.ProductService
	h.Logger = log.New(VerboseWriter(&h.LogOutput), "", log.LstdFlags)
	return h
}

func TestProductService_Product(t *testing.T) {
	t.Run("OK", testProductService_Product)
	t.Run("NotFound", testProductService_Product_NotFound)
	t.Run("ErrInternal", testProductService_Product_ErrInternal)
}

func testProductService_Product(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductFn = func(id fruit.ProductID) (*fruit.Product, error) {
		return &fruit.Product{ID: "A"}, nil
	}

	// Retrieve product.
	p, err := c.ProductService().Product(fruit.ProductID("A"))
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(p, &fruit.Product{ID: "A"}) {
		t.Fatalf("unexpected product: %+v", p)
	}
}

func testProductService_Product_NotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductFn = func(id fruit.ProductID) (*fruit.Product, error) {
		return nil, nil
	}

	// Retrieve product.
	if d, err := c.ProductService().Product(fruit.ProductID("NO SUCH PRODUCT")); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Fatal("unexpected nil product")
	}
}

func testProductService_Product_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductFn = func(id fruit.ProductID) (*fruit.Product, error) {
		return nil, errors.New("marker")
	}

	// Retrieve product.
	if p, err := c.ProductService().Product(fruit.ProductID("XXX")); err != fruit.ErrInternal {
		t.Fatal(err)
	} else if p != nil {
		t.Fatal("unexpected nil product")
	}
}

func TestProductService_Products(t *testing.T) {
	t.Run("OK", testProductService_Products)
	t.Run("NotFound", testProductService_Products_NotFound)
	t.Run("ErrInternal", testProductService_Products_ErrInternal)
}

func testProductService_Products(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductsFn = func() ([]*fruit.Product, error) {
		var products []*fruit.Product
		products = append(products, &fruit.Product{ID: "A"}, &fruit.Product{ID: "B"})

		return products, nil
	}

	if p, err := c.ProductService().Products(); err != nil {
		t.Fatal(err)
	} else if len(p) != 2 {
		t.Fatalf("expected to return two products but returned: %+v", p)
	}
}

func testProductService_Products_NotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductsFn = func() ([]*fruit.Product, error) {
		return nil, nil
	}

	// Retrieve products.
	if d, err := c.ProductService().Products(); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Fatal("unexpected nil product")
	}
}

func testProductService_Products_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductsFn = func() ([]*fruit.Product, error) {
		return nil, errors.New("marker")
	}

	// Retrieve product.
	if p, err := c.ProductService().Products(); err != fruit.ErrInternal {
		t.Fatal(err)
	} else if p != nil {
		t.Fatal("unexpected nil product")
	}
}

func TestProductService_Create(t *testing.T) {
	t.Run("OK", testProductService_CreateProduct)
	t.Run("ErrProductRequired", testProductService_CreateProduct_ErrProductRequired)
	t.Run("ErrProductExists", testProductService_CreateProduct_ErrProductExists)
	t.Run("ErrProductIDRequired", testProductService_CreateProduct_ErrProductIDRequired)
	t.Run("ErrInternal", testProductService_Products_ErrInternal)
}

func testProductService_CreateProduct(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *fruit.Product) error {
		if !reflect.DeepEqual(p, &fruit.Product{ID: "XXX", Token: "TOKEN"}) {
			t.Fatalf("unexpected product: %v", p)
		}

		// Update mod time.
		p.ModTime = Now

		return nil
	}

	p := &fruit.Product{ID: "XXX", Token: "TOKEN"}

	// Create product.
	err := c.ProductService().CreateProduct(p)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(p, &fruit.Product{ID: "XXX", Token: "TOKEN", ModTime: Now}) {
		t.Fatalf("unexpected product: %v", p)
	}
}

func testProductService_CreateProduct_ErrProductRequired(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *fruit.Product) error {
		return fruit.ErrProductRequired
	}

	if err := c.ProductService().CreateProduct(nil); err != fruit.ErrProductRequired {
		t.Fatal(err)
	}
}

func testProductService_CreateProduct_ErrProductExists(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *fruit.Product) error {
		return fruit.ErrProductExists
	}

	if err := c.ProductService().CreateProduct(&fruit.Product{ID: "XXX"}); err != fruit.ErrProductExists {
		t.Fatal(err)
	}
}

func testProductService_CreateProduct_ErrProductIDRequired(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *fruit.Product) error {
		return fruit.ErrProductIDRequired
	}

	if err := c.ProductService().CreateProduct(&fruit.Product{}); err != fruit.ErrProductIDRequired {
		t.Fatal(err)
	}
}
func testProductService_CreateProduct_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *fruit.Product) error {
		return errors.New("marker")
	}

	if err := c.ProductService().CreateProduct(&fruit.Product{}); err != fruit.ErrInternal {
		t.Fatal(err)
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	t.Run("OK", testProductService_UpdateProduct)
	t.Run("NotFound", testProductService_UpdateProduct_ErrProductNotFound)
	t.Run("ErrInternal", testProductService_UpdateProduct_ErrInternal)
}

func testProductService_UpdateProduct(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.UpdateProductFn = func(id fruit.ProductID, p *fruit.Product) error {
		// Update mod time.
		p.ModTime = Now
		p.ID = id

		return nil
	}

	p := &fruit.Product{Token: "TOKEN"}

	// Update product.
	err := c.ProductService().UpdateProduct(fruit.ProductID("XXX"), p)
	if err != nil {
		t.Fatal(err)
	} else if p.ID != "XXX" {
		t.Fatalf("product failed to update: %v", p)
	}
}

func testProductService_UpdateProduct_ErrProductNotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.UpdateProductFn = func(id fruit.ProductID, p *fruit.Product) error {
		return fruit.ErrProductNotFound
	}

	// Update product.
	err := c.ProductService().UpdateProduct("XXX", &fruit.Product{ID: "XXX"})
	if err != fruit.ErrProductNotFound {
		t.Fatal(err)
	}
}

func testProductService_UpdateProduct_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.UpdateProductFn = func(id fruit.ProductID, p *fruit.Product) error {
		return errors.New("marker")
	}

	// Update product.
	err := c.ProductService().UpdateProduct("XXX", &fruit.Product{ID: "XXX"})
	if err != fruit.ErrInternal {
		t.Fatal(err)
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	// TODO: Add Token unauthorization
	t.Run("OK", testProductService_DeleteProduct)
	t.Run("ErrProductNotFound", testProductService_DeleteProduct_ErrProductNotFound)
	t.Run("ErrInternal", testProductService_DeleteProduct_ErrInternal)
}

func testProductService_DeleteProduct(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.DeleteProductFn = func(id fruit.ProductID, token string) error {
		return nil
	}

	// Delete product.
	err := c.ProductService().DeleteProduct("XXX", "TOKEN")
	if err != nil {
		t.Fatal(err)
	}
}

func testProductService_DeleteProduct_ErrProductNotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.DeleteProductFn = func(id fruit.ProductID, token string) error {
		return fruit.ErrProductNotFound
	}

	// Delete product.
	err := c.ProductService().DeleteProduct("XXX", "TOKEN")
	if err != fruit.ErrProductNotFound {
		t.Fatal(err)
	}
}

func testProductService_DeleteProduct_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.DeleteProductFn = func(id fruit.ProductID, token string) error {
		return errors.New("marker")
	}

	// Delete product.
	err := c.ProductService().DeleteProduct("XXX", "TOKEN")
	if err != fruit.ErrInternal {
		t.Fatal(err)
	}
}
