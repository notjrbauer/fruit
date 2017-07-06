package bolt_test

import (
	"errors"
	"reflect"
	"testing"

	fruit "github.com/notjrbauer/fruitvendor"
)

func TestCreateUser(t *testing.T) {
	t.Run("OK", testUserService_CreateUser)
	t.Run("ErrUserIDRequired", testUserService_CreateUser_ErrUserIDRequired)
	t.Run("ErrUserExists", testUserService_CreateUser_ErrUserExists)
}

func testUserService_CreateUser(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	s := c.UserService()

	user := fruit.User{
		ID:      "ID",
		Name:    "NAME",
		Address: &fruit.Address{},
		CardID:  "CARDID",
	}

	if err := s.CreateUser(&user); err != nil {
		t.Fatal(err)
	}

	other, err := s.User("ID")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&user, other) {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func testUserService_CreateUser_ErrUserIDRequired(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	s := c.UserService()

	user := fruit.User{
		Name:    "NAME",
		Address: &fruit.Address{},
		CardID:  "CARDID",
	}

	if err := s.CreateUser(&user); err != fruit.ErrUserIDRequired {
		t.Fatal(err)
	}
}

func testUserService_CreateUser_ErrUserExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	s := c.UserService()

	user := fruit.User{
		ID:      "ID",
		Name:    "NAME",
		Address: &fruit.Address{},
		CardID:  "CARDID",
	}

	if err := s.CreateUser(&user); err != nil {
		t.Fatal(err)
	}

	// Create same user.
	if err := s.CreateUser(&user); err != fruit.ErrUserExists {
		t.Fatal(errors.New("expected error when creating duplicate user"))
	}
}
