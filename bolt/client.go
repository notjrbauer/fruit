package bolt

import (
	"time"

	"github.com/asdine/storm"
	"github.com/notjrbauer/fruit"
)

// Client represents a client to the underlying bolt db structure.
type Client struct {
	// Filename to the BoltDB database.
	Path string

	// Returns the current time.
	Now func() time.Time

	// Services
	productService ProductService
	userService    UserService

	db *storm.DB
}

func NewClient() *Client {
	c := &Client{Now: time.Now}
	c.productService.client = c
	c.userService.client = c
	return c
}

func (c *Client) Open() error {
	db, err := storm.Open(c.Path)

	if err != nil {
		return err
	}

	c.db = db

	return nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) ProductService() fruit.ProductService {
	return &c.productService
}

func (c *Client) UserService() fruit.UserService {
	return &c.userService
}
