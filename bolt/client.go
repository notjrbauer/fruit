package bolt

import (
	"time"

	"github.com/asdine/storm"
	"github.com/notjrbauer/fruitvendor"
)

// Client represents a client to the underlying bolt db structure.
type Client struct {
	// Filename to the BoltDB database.
	Path string

	// Returns the current time.
	Now func() time.Time

	// Services
	productService ProductService

	db *storm.DB
}

func NewClient() *Client {
	c := &Client{Now: time.Now}
	c.productService.client = c
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

func (c *Client) ProductService() fruitvendor.ProductService {
	return &c.productService
}
