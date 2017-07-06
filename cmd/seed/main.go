package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/notjrbauer/fruit"
	"github.com/notjrbauer/fruit/bolt"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Parse command line arguments
	generateCommand := flag.NewFlagSet("generate", flag.ContinueOnError)
	srcDBPath := generateCommand.String("src-db", "", "source db path")
	since := generateCommand.Int("start-txtid", 0, "replay from txid")

	// First argument specifies a subcommand to run.
	switch os.Args[1] {
	case "generate":
		generateCommand.Parse(os.Args[2:])
		fmt.Fprintln(os.Stdout, srcDBPath, since)
		err := generate(*srcDBPath)
		if err != nil {
			panic(err)
		}
	}
}

func generate(path string) error {
	//for {
	// Generate temporary path.
	_, err := newFile(path)

	err = generateFile(path, 10)

	return err
}

func generateFile(path string, n int) error {
	// Initialize client.
	c := bolt.NewClient()
	c.Path = path

	if err := c.Open(); err != nil {
		return err
	}

	colors := []string{"Red", "Green", "White", "Blue"}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	for count := 0; count < n; count++ {
		color := r.Intn(len(colors))
		id := strconv.Itoa(r.Intn(1000) + count)

		// TODO: Break these into their own functions when all services are defined.
		// Generate products.
		if err := c.ProductService().CreateProduct(&fruit.Product{ID: fruit.ProductID(id), Color: colors[color]}); err != nil {
			return err
		}

		// Generate Users
		if err := c.UserService().CreateUser(&fruit.User{ID: fruit.UserID(id), Name: colors[color], CardID: strconv.Itoa(rand.Int())}); err != nil {
			return err
		}
	}

	defer c.Close()

	return nil
}

func newFile(path string) (*os.File, error) {
	// detect if file exists
	_, err := os.Stat(path)
	if os.IsExist(err) {
		return nil, err
	}
	// create file if not exists
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
