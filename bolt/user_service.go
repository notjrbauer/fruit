package bolt

import (
	fruit "github.com/notjrbauer/fruitvendor"
)

type UserService struct {
	client *Client
}

// User returns a user by ID.
func (s *UserService) User(id fruit.UserID) (*fruit.User, error) {
	// Start read-only transaction.
	tx, err := s.client.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal user.
	var u fruit.User
	bucket := s.client.db.From("Users")

	if err := bucket.One("ID", id, &u); err != nil {
		return nil, err
	} else if &u == nil {
		return nil, nil
	}

	return &u, nil
}

// Users returns a list of users.
// TODO: Add params
func (s *UserService) Users() ([]*fruit.User, error) {
	panic("not implemented")
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(u *fruit.User) error {
	// Require id
	// TODO: Don't require ID, have the DB generate it
	if u.ID == "" {
		return fruit.ErrUserIDRequired
	}

	bucket := s.client.db.From("Users")

	// Start the read-write transaction.
	tx, err := bucket.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify user doesn't already exist.
	var user fruit.User
	tx.One("ID", u.ID, &user)

	if user.ID != "" {
		return fruit.ErrUserExists
	}

	// Save the user.
	if err := tx.Save(u); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteUser removes an existing user.
func (s *UserService) DeleteUser(id fruit.UserID) error {
	panic("not implemented")
}

// UpdateUser removes an existing user.
func (s *UserService) UpdateUser(id fruit.UserID, u *fruit.User) error {
	panic("not implemented")
}
