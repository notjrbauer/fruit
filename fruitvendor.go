package fruit

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
// TODO: Decide if we really need to use client and not
// just standalone services.
type Client interface {
	ProductService() ProductService
	UserService() UserService
}

// ProductService represents a service for managing products
type ProductService interface {
	Product(id ProductID) (*Product, error)
	Products() ([]*Product, error)
	CreateProduct(p *Product) error
	UpdateProduct(id ProductID, p *Product) error
	DeleteProduct(id ProductID, token string) error
}

type Address struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

type UserID string

type User struct {
	ID      UserID   `json:"userID" storm:"id"`
	Name    string   `json:"name"`
	Address *Address `json:"address"`
	CardID  string   `json:"card"`
}

type UserService interface {
	User(id UserID) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	DeleteUser(id UserID) error
	UpdateUser(id UserID, u *User) error
}

type TransactionID string

type Transaction struct {
	ID     TransactionID `json:"transactionID" validate:"nonzero"`
	UserID UserID        `json:"userID"`
	Count  int           `json:"count"`
	Active bool          `json:"active"`
}

type TransactionService interface {
	Transaction(id TransactionID) (*Transaction, error)
	Transactions(id UserID) ([]*Transaction, error)
	CreateTransaction(t *Transaction) error
	UpdateTransaction(id TransactionID, t *Transaction) error
	DeleteTransaction(id TransactionID) error
}
