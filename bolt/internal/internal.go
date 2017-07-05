package internal

import (
	"encoding/json"

	"github.com/notjrbauer/fruitvendor"
)

func MarshalProduct(p *fruit.Product) ([]byte, error) {
	return json.Marshal(p)
}

func UnmarshalProduct(b []byte, p *fruit.Product) error {
	return json.Unmarshal(b, p)
}
