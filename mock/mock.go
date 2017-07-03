package mock

import (
	"github.com/notjrbauer/caps"
)

type ProductService struct {
	ProductFn      func(id caps.ProductID) (*caps.Product, error)
	ProductInvoked bool

	ProductsFn      func() ([]*caps.Product, error)
	ProductsInvoked bool

	CreateProductFn      func(p *caps.Product) error
	CreateProductInvoked bool

	UpdateProductFn      func(id caps.ProductID, p *caps.Product) error
	UpdateProductInvoked bool

	DeleteProductFn      func(id caps.ProductID, token string) error
	DeleteProductInvoked bool
}

func (s *ProductService) Product(id caps.ProductID) (*caps.Product, error) {
	s.ProductInvoked = true
	return s.ProductFn(id)
}

func (s *ProductService) Products() ([]*caps.Product, error) {
	s.ProductsInvoked = true
	return s.ProductsFn()
}

func (s *ProductService) CreateProduct(p *caps.Product) error {
	s.CreateProductInvoked = true
	return s.CreateProductFn(p)
}

func (s *ProductService) UpdateProduct(id caps.ProductID, p *caps.Product) error {
	s.UpdateProductInvoked = true
	return s.UpdateProductFn(id, p)
}

func (s *ProductService) DeleteProduct(id caps.ProductID, token string) error {
	s.DeleteProductInvoked = true
	return s.DeleteProductFn(id, token)
}
