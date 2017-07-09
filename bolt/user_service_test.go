package bolt_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/notjrbauer/fruit"
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

func TestDeleteUser(t *testing.T) {
	t.Run("OK", testUserService_DeleteUser)
}

func testUserService_DeleteUser(t *testing.T) {
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
	if err := s.DeleteUser(user.ID); err != nil {
		t.Fatal(err)
	}

	// User should not exist
	if _, err := s.User(user.ID); err == nil {
		t.Fatal(errors.New("expected error when removing user"))
	}
}

func TestUpdateUser(t *testing.T) {
	t.Run("OK", testUserService_UpdateUser)
}

func testUserService_UpdateUser(t *testing.T) {
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

	user.Name = "UPDATE_NAME"
	user.CardID = "UPDATE_CARDID"
	if err := s.UpdateUser(user.ID, &user); err != nil {
		t.Fatal(err)
	}

	// User should be updated.
	if u, err := s.User(user.ID); err != nil {
		t.Fatal(errors.New("expected error when updating user"))
	} else if u.Name != "UPDATE_NAME" {
		t.Fatalf("unexpected product sku: %s", u.Name)
	}
}
