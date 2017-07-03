package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/notjrbauer/fruitvendor/bolt"
	"github.com/notjrbauer/fruitvendor/http"
)

func main() {
	c := bolt.NewClient()
	c.Path = "../seed/boltdbseed.db"
	err := c.Open()
	if err != nil {
		panic(err)
	}

	s := http.NewServer()
	s.Handler = &http.Handler{
		ProductHandler: http.NewProductHandler(),
	}
	s.Handler.ProductHandler.ProductService = c.ProductService()
	s.Addr = ":3000"
	_ = s.Open()
	spew.Dump(s)

}
