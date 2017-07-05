package internal_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/notjrbauer/fruitvendor"
	"github.com/notjrbauer/fruitvendor/bolt/internal"
)

func TestMarshalProduct(t *testing.T) {
	v := fruit.Product{
		ID:          "ID",
		SKU:         "SKU",
		Name:        "NAME",
		Type:        "TYPE",
		Color:       "COLOR",
		Description: "DESCRIPTION",
		ModTime:     time.Now().UTC(),
	}

	var other fruit.Product
	if buf, err := internal.MarshalProduct(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.UnmarshalProduct(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
