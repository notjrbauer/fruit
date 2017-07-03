package bolt_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/notjrbauer/fruitvendor/bolt"
)

var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

type Client struct {
	*bolt.Client
}

func NewClient() *Client {

	f, err := ioutil.TempFile("", "fruitvendor-bolt-client-")
	if err != nil {
		panic(err)
	}

	f.Close()

	c := &Client{
		Client: bolt.NewClient(),
	}
	c.Path = f.Name()
	c.Now = func() time.Time { return Now }

	return c
}

// MustOpenClient returns a new, open instance of Client
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

func (c *Client) Close() error {
	defer os.Remove(c.Path)
	return c.Client.Close()
}
