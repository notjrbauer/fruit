package internal

import (
	"encoding/json"

	"github.com/notjrbauer/fruitvendor"
)

func MarshalProduct(p *fruitvendor.Product) ([]byte, error) {
	return json.Marshal(p)
}

func UnmarshalProduct(b []byte, p *fruitvendor.Product) error {
	return json.Unmarshal(b, p)
}
