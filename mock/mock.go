package mock

import (
	"github.com/notjrbauer/fruitvendor"
)

type ProductService struct {
	ProductFn      func(id fruitvendor.ProductID) (*fruitvendor.Product, error)
	ProductInvoked bool

	ProductsFn      func() ([]*fruitvendor.Product, error)
	ProductsInvoked bool

	CreateProductFn      func(p *fruitvendor.Product) error
	CreateProductInvoked bool

	UpdateProductFn      func(id fruitvendor.ProductID, p *fruitvendor.Product) error
	UpdateProductInvoked bool

	DeleteProductFn      func(id fruitvendor.ProductID, token string) error
	DeleteProductInvoked bool
}

func (s *ProductService) Product(id fruitvendor.ProductID) (*fruitvendor.Product, error) {
	s.ProductInvoked = true
	return s.ProductFn(id)
}

func (s *ProductService) Products() ([]*fruitvendor.Product, error) {
	s.ProductsInvoked = true
	return s.ProductsFn()
}

func (s *ProductService) CreateProduct(p *fruitvendor.Product) error {
	s.CreateProductInvoked = true
	return s.CreateProductFn(p)
}

func (s *ProductService) UpdateProduct(id fruitvendor.ProductID, p *fruitvendor.Product) error {
	s.UpdateProductInvoked = true
	return s.UpdateProductFn(id, p)
}

func (s *ProductService) DeleteProduct(id fruitvendor.ProductID, token string) error {
	s.DeleteProductInvoked = true
	return s.DeleteProductFn(id, token)
}
