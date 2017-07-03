package http_test

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/notjrbauer/caps"
	"github.com/notjrbauer/caps/http"
	"github.com/notjrbauer/caps/mock"
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
	s.Handler.ProductHandler.ProductService.ProductFn = func(id caps.ProductID) (*caps.Product, error) {
		return &caps.Product{ID: "A"}, nil
	}

	// Retrieve product.
	p, err := c.ProductService().Product(caps.ProductID("A"))
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(p, &caps.Product{ID: "A"}) {
		t.Fatalf("unexpected product: %+v", p)
	}
}

func testProductService_Product_NotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductFn = func(id caps.ProductID) (*caps.Product, error) {
		return nil, nil
	}

	// Retrieve product.
	if d, err := c.ProductService().Product(caps.ProductID("NO SUCH PRODUCT")); err != nil {
		t.Fatal(err)
	} else if d != nil {
		t.Fatal("unexpected nil product")
	}
}

func testProductService_Product_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.ProductHandler.ProductService.ProductFn = func(id caps.ProductID) (*caps.Product, error) {
		return nil, errors.New("marker")
	}

	// Retrieve product.
	if p, err := c.ProductService().Product(caps.ProductID("XXX")); err != caps.ErrInternal {
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
	s.Handler.ProductHandler.ProductService.ProductsFn = func() ([]*caps.Product, error) {
		var products []*caps.Product
		products = append(products, &caps.Product{ID: "A"}, &caps.Product{ID: "B"})

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
	s.Handler.ProductHandler.ProductService.ProductsFn = func() ([]*caps.Product, error) {
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
	s.Handler.ProductHandler.ProductService.ProductsFn = func() ([]*caps.Product, error) {
		return nil, errors.New("marker")
	}

	// Retrieve product.
	if p, err := c.ProductService().Products(); err != caps.ErrInternal {
		t.Fatal(err)
	} else if p != nil {
		t.Fatal("unexpected nil product")
	}
}

func TestProductService_Create(t *testing.T) {
	t.Run("OK", testProductService_CreateProduct)
	t.Run("ErrProductRequired", testProductService_CreateProduct_ErrProductRequired)
	t.Run("ErrProductExists", testProductService_CreateProduct_ErrProductExists)
	t.Run("ErrProductExists", testProductService_CreateProduct_ErrProductIDRequired)
	//t.Run("NotFound", testProductService_Products_NotFound)
	//t.Run("ErrInternal", testProductService_Products_ErrInternal)
}

func testProductService_CreateProduct(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock server.
	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *caps.Product) error {
		if !reflect.DeepEqual(p, &caps.Product{ID: "XXX", Token: "TOKEN"}) {
			t.Fatalf("unexpected product: %v", p)
		}

		// Update mod time.
		p.ModTime = Now

		return nil
	}

	p := &caps.Product{ID: "XXX", Token: "TOKEN"}

	// Create product.
	err := c.ProductService().CreateProduct(p)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(p, &caps.Product{ID: "XXX", Token: "TOKEN", ModTime: Now}) {
		t.Fatalf("unexpected product: %v", p)
	}
}

func testProductService_CreateProduct_ErrProductRequired(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *caps.Product) error {
		return caps.ErrProductRequired
	}

	if err := c.ProductService().CreateProduct(nil); err != caps.ErrProductRequired {
		t.Fatal(err)
	}
}

func testProductService_CreateProduct_ErrProductExists(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *caps.Product) error {
		return caps.ErrProductExists
	}

	if err := c.ProductService().CreateProduct(&caps.Product{ID: "XXX"}); err != caps.ErrProductExists {
		t.Fatal(err)
	}
}

func testProductService_CreateProduct_ErrProductIDRequired(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *caps.Product) error {
		return caps.ErrProductIDRequired
	}

	if err := c.ProductService().CreateProduct(&caps.Product{}); err != caps.ErrProductIDRequired {
		t.Fatal(err)
	}
}
func testProductService_CreateProduct_ErrInternal(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.ProductHandler.ProductService.CreateProductFn = func(p *caps.Product) error {
		return errors.New("marker")
	}

	if err := c.ProductService().CreateProduct(&caps.Product{}); err != caps.ErrInternal {
		t.Fatal(err)
	}
}
