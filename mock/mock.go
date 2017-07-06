package mock

import (
	"github.com/notjrbauer/fruit"
)

type ProductService struct {
	ProductFn      func(id fruit.ProductID) (*fruit.Product, error)
	ProductInvoked bool

	ProductsFn      func() ([]*fruit.Product, error)
	ProductsInvoked bool

	CreateProductFn      func(p *fruit.Product) error
	CreateProductInvoked bool

	UpdateProductFn      func(id fruit.ProductID, p *fruit.Product) error
	UpdateProductInvoked bool

	DeleteProductFn      func(id fruit.ProductID, token string) error
	DeleteProductInvoked bool
}

func (s *ProductService) Product(id fruit.ProductID) (*fruit.Product, error) {
	s.ProductInvoked = true
	return s.ProductFn(id)
}

func (s *ProductService) Products() ([]*fruit.Product, error) {
	s.ProductsInvoked = true
	return s.ProductsFn()
}

func (s *ProductService) CreateProduct(p *fruit.Product) error {
	s.CreateProductInvoked = true
	return s.CreateProductFn(p)
}

func (s *ProductService) UpdateProduct(id fruit.ProductID, p *fruit.Product) error {
	s.UpdateProductInvoked = true
	return s.UpdateProductFn(id, p)
}

func (s *ProductService) DeleteProduct(id fruit.ProductID, token string) error {
	s.DeleteProductInvoked = true
	return s.DeleteProductFn(id, token)
}
